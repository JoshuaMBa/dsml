package server_lib

import (
	"context"
	"errors"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/JoshuaMBa/dsml/failure_injection"
	fipb "github.com/JoshuaMBa/dsml/failure_injection/proto"
	"github.com/JoshuaMBa/dsml/gpu_sim/proto"
	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"

	// "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// options for service init
type GPUDeviceOptions struct {
	// failure injection config params
	SleepNs              int64
	FailureRate          int64
	ResponseOmissionRate int64

	////////////////////////////
	// Stuff added by michael //
	////////////////////////////
	DeviceId     uint64
	MemoryNeeded uint64 // size of memory needed in bytes

}

func DefaultGPUDeviceOptions() *GPUDeviceOptions {
	return &GPUDeviceOptions{
		SleepNs:              0,
		FailureRate:          0,
		ResponseOmissionRate: 0,
		DeviceId:             11,
		MemoryNeeded:         0x1000,
	}
}

type GPUDeviceServer struct {
	proto.UnimplementedGPUDeviceServer
	// options are read-only and intended to be immutable during the lifetime of
	// the service.  The failure injection config (in the FailureInjector) is
	// mutable (via SetInjectionConfigRequest-s) after the server starts.
	options GPUDeviceOptions
	fi      *failure_injection.FailureInjector

	////////////////////////////
	// Stuff added by michael //
	////////////////////////////

	////////////////////////////
	//    simulate memory     //
	////////////////////////////

	memory map[uint64][]byte // gpu's memory space
	mu     sync.Mutex        // thread safe memory

	////////////////////////////
	// device info (i don't anticipate ever using this, maybe it goes into options?)
	////////////////////////////
	deviceId   uint64
	minMemAddr uint64
	maxMemAddr uint64

	////////////////////////////
	// gpu communication logical
	////////////////////////////
	rank   uint32 // my rank in the communicator
	nRanks uint32 // total number of gpus in the communicator

	rankToAddress map[uint32]string // map between rank and addresses

	////////////////////////////
	// gpu communications, implementation detail
	////////////////////////////
	streamId atomic.Uint64         // my streamId when sending to others
	peers    []*pb.GPUDeviceClient // rpc handles for other gpus

	streamSrc      map[uint64]StreamSrcInfo // where the streams i send come from: (lookup is streamID->{src addr + size})
	streamDest     map[uint64]uint64        // where streams i am receiving should go (lookup is rank->streamId->memAddr (combine into one struct, streamHandler?)
	currentStream  map[uint32]uint64        // which stream for each sender (lookup is rank->streamId)
	streamStatuses sync.Map
}

type StreamSrcInfo struct {
	SendBuffAddr uint64
	NumBytes     uint64
	DstRank      uint32
}

func MakeGPUDeviceServer(
	options GPUDeviceOptions,
) (*GPUDeviceServer, error) {
	fi := failure_injection.MakeFailureInjector()
	fi.SetInjectionConfig(
		options.SleepNs,
		options.FailureRate,
		options.ResponseOmissionRate,
	)
	log.Printf(
		"Starting GPUDevice with failure injection config: [sleepNs: %d, failureRate: %d, responseOmissionRate: %d]",
		options.SleepNs,
		options.FailureRate,
		options.ResponseOmissionRate,
	)
	return &GPUDeviceServer{
		options:        options,
		fi:             fi,
		deviceId:       options.DeviceId,
		minMemAddr:     0,
		maxMemAddr:     options.MemoryNeeded,
		memory:         make(map[uint64][]byte),
		streamSrc:      make(map[uint64]StreamSrcInfo),
		streamDest:     make(map[uint64]uint64),
		currentStream:  make(map[uint32]uint64),
		streamStatuses: sync.Map{},
	}, nil
}

func MockGPUDeviceServer(deviceId, minMemAddr, maxMemAddr uint64, rankToAddress map[uint32]string) *GPUDeviceServer {
	return &GPUDeviceServer{
		deviceId:      deviceId,
		minMemAddr:    minMemAddr,
		maxMemAddr:    maxMemAddr,
		rankToAddress: rankToAddress,
	}
}

func (gpu *GPUDeviceServer) SetupCommunication(
	ctx context.Context,
	req *pb.SetupCommunicationRequest,
) (*pb.SetupCommunicationResponse, error) {
	gpu.mu.Lock()
	defer gpu.mu.Unlock()
	if req == nil || req.RankToAddress == nil {
		return nil, errors.New("invalid SetupCommunicationRequest: rankToAddress is nil")
	}
	if _, exists := req.RankToAddress[req.Rank.Value]; !exists {
		return nil, errors.New("current rank not found in rankToAddress")
	}

	gpu.rank = req.Rank.Value
	gpu.rankToAddress = req.RankToAddress
	gpu.nRanks = uint32(len(req.RankToAddress))

	return &pb.SetupCommunicationResponse{
		Success: true,
	}, nil
}

func (gpu *GPUDeviceServer) GetDeviceMetadata(
	ctx context.Context,
	req *pb.GetDeviceMetadataRequest,
) (*pb.GetDeviceMetadataResponse, error) {
	return &pb.GetDeviceMetadataResponse{
		Metadata: &pb.DeviceMetadata{
			DeviceId:   &pb.DeviceId{Value: gpu.deviceId},  // Wrap the DeviceId
			MinMemAddr: &pb.MemAddr{Value: gpu.minMemAddr}, // Wrap MinMemAddr
			MaxMemAddr: &pb.MemAddr{Value: gpu.maxMemAddr}, // Wrap MaxMemAddr
		},
	}, nil
}

func (gpu *GPUDeviceServer) BeginSend(
	ctx context.Context,
	req *pb.BeginSendRequest,
) (*pb.BeginSendResponse, error) {
	// get current stream id and increment
	// kick off goroutine actually sending data using streamsend
	streamId := gpu.streamId.Load()

	gpu.streamSrc[streamId-1] = StreamSrcInfo{req.SendBuffAddr.Value, req.NumBytes, req.DstRank.Value}
	gpu.streamStatuses.Store(streamId, pb.Status_READY)

	gpu.streamId.Add(1)

	// add error checking

	return &pb.BeginSendResponse{Initiated: true, StreamId: &pb.StreamId{Value: streamId}}, nil
}

func (gpu *GPUDeviceServer) BeginReceive(
	ctx context.Context,
	req *pb.BeginReceiveRequest,
) (*pb.BeginReceiveResponse, error) {
	gpu.mu.Lock()
	gpu.streamDest[req.StreamId.Value] = req.RecvBuffAddr.Value
	gpu.mu.Unlock()

	gpu.streamStatuses.Store(req.StreamId.Value, pb.Status_IN_PROGRESS)

	return &pb.BeginReceiveResponse{
		Initiated: true,
	}, nil
}

func (gpu *GPUDeviceServer) StreamSend(
	stream pb.GPUDevice_StreamSendServer,
) error {
	var streamId uint64

	for {
		// Receive a DataChunk message from the stream
		req, err := stream.Recv()
		if err == io.EOF {
			gpu.streamStatuses.Store(streamId, pb.Status_SUCCESS)

			// End of stream, send response
			response := &pb.StreamSendResponse{
				Success: true,
			}
			return stream.SendAndClose(response)
		}
		if err != nil {
			// Handle errors during streaming
			gpu.streamStatuses.Store(streamId, pb.Status_FAILED)
			return status.Errorf(codes.Internal, "failed to receive stream: %v", err)
		}

		// Simulate writing the data to memory
		gpu.mu.Lock()
		for addr, data := range req.Data {
			uintAddr := uint64(addr)

			if uintAddr < gpu.minMemAddr || uintAddr >= gpu.maxMemAddr {
				gpu.mu.Unlock()
				return status.Errorf(codes.InvalidArgument, "invalid memory address: %d", uintAddr)
			}

			gpu.memory[uintAddr] = []byte{data}
		}
		gpu.mu.Unlock()
	}
}

func (gpu *GPUDeviceServer) GetStreamStatus(
	ctx context.Context,
	req *pb.GetStreamStatusRequest,
) (*pb.GetStreamStatusResponse, error) {
	streamStatus, ok := gpu.streamStatuses.Load(req.StreamId.Value)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "stream ID %d not found", req.StreamId.Value)
	}

	return &pb.GetStreamStatusResponse{
		Status: streamStatus.(pb.Status),
	}, nil
}

func (gpu *GPUDeviceServer) SetInjectionConfig(
	ctx context.Context,
	req *fipb.SetInjectionConfigRequest,
) (*fipb.SetInjectionConfigResponse, error) {
	gpu.fi.SetInjectionConfigPb(req.Config)
	log.Printf(
		"GPUDevice failure injection config set to: [sleepNs: %d, failureRate: %d, responseOmissionRate: %d]",
		req.Config.SleepNs,
		req.Config.FailureRate,
		req.Config.ResponseOmissionRate,
	)

	return &fipb.SetInjectionConfigResponse{}, nil
}

func (gpu *GPUDeviceServer) StreamSendThread() {
	for {
		// calls streamsend on things in streamstatus which are not ready!
		// for now will just acquire lock
		gpu.streamStatuses.Range(func(streamId, status interface{}) bool {
			if status == pb.Status_READY {
				return false
			} else {
				return true // true to continue
			}
		})
	}
}
