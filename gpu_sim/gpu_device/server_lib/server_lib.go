package server_lib

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/JoshuaMBa/dsml/failure_injection"
	fipb "github.com/JoshuaMBa/dsml/failure_injection/proto"
	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"

	// "google.golang.org/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// options for service init
type GPUDeviceOptions struct {
	// Failure injection config params
	SleepNs              int64  `json:"sleep_ns"`
	FailureRate          int64  `json:"failure_rate"`
	ResponseOmissionRate int64  `json:"response_omission_rate"`
	DeviceId             uint64 `json:"device_id"`
	MemoryNeeded         uint64 `json:"memory_needed"` // Size of memory needed in bytes
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
	pb.UnimplementedGPUDeviceServer
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

	memory []byte
	mu     sync.Mutex // thread safe memory

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

	rank2Address map[uint32]string // map between rank and addresses

	////////////////////////////
	// gpu communications, implementation detail
	////////////////////////////
	streamId atomic.Uint64 // my streamId when sending to others

	conns []*grpc.ClientConn
	peers []pb.GPUDeviceClient // rpc handles for other gpus

	streamSrc chan StreamSrcInfo

	streamDest     map[uint32]map[uint64]uint64 // where streams i am receiving should go (lookup is rank->streamId->memAddr (combine into one struct, streamHandler?)
	currentStream  map[uint32]uint64            // which stream for each sender (lookup is rank->streamId)
	streamStatuses sync.Map

	streamSendSignal chan int
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

	// preallocate memory, if it's big it should use mmap
	mem := make([]byte, options.MemoryNeeded)
	startAddr := uint64(uintptr(unsafe.Pointer(&mem[0])))
	endAddr := startAddr + options.MemoryNeeded

	return &GPUDeviceServer{
		options:          options,
		fi:               fi,
		deviceId:         options.DeviceId,
		minMemAddr:       startAddr,
		maxMemAddr:       endAddr,
		memory:           mem,
		streamSrc:        make(chan StreamSrcInfo, 16),
		streamDest:       make(map[uint32]map[uint64]uint64),
		currentStream:    make(map[uint32]uint64),
		streamStatuses:   sync.Map{},
		streamSendSignal: make(chan int),
	}, nil
}

func MockGPUDeviceServer(deviceId, minMemAddr, maxMemAddr uint64, rankToAddress map[uint32]string) *GPUDeviceServer {
	return &GPUDeviceServer{
		deviceId:     deviceId,
		minMemAddr:   minMemAddr,
		maxMemAddr:   maxMemAddr,
		rank2Address: rankToAddress,
	}
}

func (gpu *GPUDeviceServer) SetupCommunication(
	ctx context.Context,
	req *pb.SetupCommunicationRequest,
) (*pb.SetupCommunicationResponse, error) {
	// you should not be allowed to call this more than once
	gpu.mu.Lock()
	defer gpu.mu.Unlock()
	if req == nil || req.RankToAddress == nil {
		return nil, status.Error(codes.InvalidArgument, "missing RankToAddress mappings")
	}
	log.Printf("rank2address is : %v", req.RankToAddress)
	log.Printf("my rank is : %v", req.Rank.Value)
	if _, exists := req.RankToAddress[req.Rank.Value]; !exists {
		return nil, status.Error(codes.InvalidArgument, "did not specify rank for device")
	}

	gpu.rank = req.Rank.Value
	gpu.rank2Address = req.RankToAddress
	gpu.nRanks = uint32(len(req.RankToAddress))

	// setup connections with peers
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	gpu.conns = make([]*grpc.ClientConn, gpu.nRanks)
	gpu.peers = make([]pb.GPUDeviceClient, gpu.nRanks)
	for rank := range gpu.nRanks {
		if rank == gpu.rank {
			continue
		}

		// handle this error!
		gpu.conns[rank], _ = grpc.NewClient(gpu.rank2Address[rank], opts...)
		gpu.peers[rank] = pb.NewGPUDeviceClient(gpu.conns[rank])
	}

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

	// TODO: Add range checks!
	start := req.SendBuffAddr.Value
	end := start + req.NumBytes
	if start < gpu.minMemAddr || end >= gpu.maxMemAddr {
		return &pb.BeginSendResponse{Initiated: false},
			status.Errorf(codes.InvalidArgument,
				"invalid argument: source address range [0x%x, 0x%x) not in GPU memory range [0x%x, 0x%x)",
				start, end, gpu.minMemAddr, gpu.maxMemAddr)
	}

	streamId := gpu.streamId.Load()

	gpu.streamSrc <- StreamSrcInfo{req.SendBuffAddr.Value, req.NumBytes, req.DstRank.Value}
	gpu.streamStatuses.Store(streamId, pb.Status_READY)

	gpu.streamId.Add(1)

	return &pb.BeginSendResponse{Initiated: true, StreamId: &pb.StreamId{Value: streamId}}, nil
}

func (gpu *GPUDeviceServer) BeginReceive(
	ctx context.Context,
	req *pb.BeginReceiveRequest,
) (*pb.BeginReceiveResponse, error) {
	gpu.mu.Lock()
	gpu.streamDest[req.SrcRank.Value][req.StreamId.Value] = req.RecvBuffAddr.Value
	gpu.mu.Unlock()

	gpu.streamStatuses.Store(req.StreamId.Value, pb.Status_IN_PROGRESS)

	return &pb.BeginReceiveResponse{
		Initiated: true,
	}, nil
}

func (gpu *GPUDeviceServer) StreamSend(
	stream pb.GPUDevice_StreamSendServer,
) error {
	var rank uint32
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
		offset := gpu.streamDest[rank][streamId] - gpu.minMemAddr
		copy(gpu.memory[offset:], req.Data)
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
	fmt.Println("started streamsend thread")
	for streamInfo := range gpu.streamSrc {
		srcAddr, numBytes, dstRank := streamInfo.SendBuffAddr, streamInfo.NumBytes, streamInfo.DstRank

		// actually put this shit in a goroutine?
		// what to do in error cases?
		rpc_client := gpu.peers[dstRank]

		stream, err := rpc_client.StreamSend(context.Background())
		if err != nil {
			log.Printf("GPUDevice failed to start stream to rank %d", dstRank)
			continue
		}

		// not streaming in chunks. in cuda, you specify a datatype, which makes this make sense
		offset := srcAddr - gpu.minMemAddr
		stream.Send(&pb.DataChunk{Data: gpu.memory[offset : offset+numBytes]})

		_, err = stream.CloseAndRecv()
		if err != nil {
			log.Printf("GPUDevice failed to receive stream response from rank %d", dstRank)
			log.Printf("received error: %v", err)
		}
	}
}

func (gpu *GPUDeviceServer) GracefulStop() {
	close(gpu.streamSrc)
	for rank, conn := range gpu.conns {
		if uint32(rank) == gpu.rank {
			continue
		}
		defer conn.Close()
	}
}
