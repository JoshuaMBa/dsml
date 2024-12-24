package server_lib

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/JoshuaMBa/dsml/failure_injection"
	fipb "github.com/JoshuaMBa/dsml/failure_injection/proto"
	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
	convert "github.com/JoshuaMBa/dsml/utils"

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

	memory MemorySpace

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

	////////////////////////////
	// gpu communications, implementation detail
	////////////////////////////
	resetLock sync.RWMutex // allows for reset

	conns []*grpc.ClientConn
	peers []pb.GPUDeviceClient // rpc handles for other gpus

	streamId  atomic.Uint64 // my streamId when sending to others
	streamSrc chan StreamSrcInfo

	streamDst    []*StreamDstMonitor // where streams i am receiving should go (lookup is rank->streamId->memAddr (combine into one struct, streamHandler?)
	streamStatus []*sync.Map         // lookup is srcRank->streamid, tells status of streams being received, for streams that this gpu is sending, lookup stats[gpu.rank]
}

type StreamSrcInfo struct {
	StreamId     uint64
	SendBuffAddr uint64
	NumBytes     uint64
	DstRank      uint32
}

type StreamDstInfo struct {
	DstAddr uint64
	op      pb.ReduceOp
}

type StreamDstMonitor struct {
	cond sync.Cond
	info map[uint64]*StreamDstInfo // maps streamid->(dstaddr, op) (for a given rank only)
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
	mem, _ := NewDeviceMemory(fmt.Sprintf("device_%d_memory", options.DeviceId))
	mem.Allocate(options.MemoryNeeded)

	gpu := &GPUDeviceServer{
		options:    options,
		fi:         fi,
		deviceId:   options.DeviceId,
		memory:     mem,
		minMemAddr: mem.MinAddr(),
		maxMemAddr: mem.MaxAddr(),
		streamSrc:  make(chan StreamSrcInfo, 16),
		// we don't do streamdst until we have all the ranks?
		// ditto for streamstatus
	}

	go gpu.StreamSendDispatch()

	return gpu, nil
}

func MockGPUDeviceServer(deviceId, minMemAddr, maxMemAddr uint64, rankToAddress map[uint32]string) *GPUDeviceServer {
	return &GPUDeviceServer{
		deviceId:   deviceId,
		minMemAddr: minMemAddr,
		maxMemAddr: maxMemAddr,
	}
}

func (gpu *GPUDeviceServer) SetupCommunication(
	ctx context.Context,
	req *pb.SetupCommunicationRequest,
) (*pb.SetupCommunicationResponse, error) {
	gpu.resetLock.RLock()
	defer gpu.resetLock.RUnlock()
	shouldError := gpu.fi.MaybeInject()
	if shouldError {
		return nil, status.Error(
			codes.Internal,
			"GPU: (injected) internal error!",
		)
	}
	// you should not be allowed to call this more than once
	if req == nil || req.RankToAddress == nil {
		return nil, status.Error(codes.InvalidArgument, "missing RankToAddress mappings")
	}

	if _, exists := req.RankToAddress[req.Rank.Value]; !exists {
		return nil, status.Error(codes.InvalidArgument, "did not specify rank for device")
	}

	gpu.rank = req.Rank.Value
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
		gpu.conns[rank], _ = grpc.NewClient(req.RankToAddress[rank], opts...)
		gpu.peers[rank] = pb.NewGPUDeviceClient(gpu.conns[rank])
	}

	// setup streamdst info
	gpu.streamDst = make([]*StreamDstMonitor, gpu.nRanks)
	for i := range gpu.streamDst {
		gpu.streamDst[i] = &StreamDstMonitor{
			cond: *sync.NewCond(&sync.Mutex{}),
			info: make(map[uint64]*StreamDstInfo),
		}
	}

	// setup statuses
	gpu.streamStatus = make([]*sync.Map, gpu.nRanks)
	for i := range gpu.streamStatus {
		gpu.streamStatus[i] = &sync.Map{}
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
	gpu.resetLock.RLock()
	defer gpu.resetLock.RUnlock()
	shouldError := gpu.fi.MaybeInject()
	if shouldError {
		return nil, status.Error(
			codes.Internal,
			"GPU: (injected) internal error!",
		)
	}
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

	// protobuf default value issue
	var dstRank uint32 = 0
	if req.DstRank != nil {
		dstRank = req.DstRank.Value
	}
	streamId := gpu.streamId.Load()

	gpu.streamSrc <- StreamSrcInfo{streamId, req.SendBuffAddr.Value, req.NumBytes, dstRank}
	gpu.streamStatus[gpu.rank].Store(streamId, pb.Status_IN_PROGRESS)

	gpu.streamId.Add(1)

	return &pb.BeginSendResponse{Initiated: true, StreamId: &pb.StreamId{Value: streamId}}, nil
}

func (gpu *GPUDeviceServer) BeginReceive(
	ctx context.Context,
	req *pb.BeginReceiveRequest,
) (*pb.BeginReceiveResponse, error) {
	gpu.resetLock.RLock()
	defer gpu.resetLock.RUnlock()
	shouldError := gpu.fi.MaybeInject()
	if shouldError {
		return nil, status.Error(
			codes.Internal,
			"GPU: (injected) internal error!",
		)
	}

	log.Printf("recv req: %v", req)

	start := req.RecvBuffAddr.Value
	end := start + req.NumBytes
	if start < gpu.minMemAddr || end >= gpu.maxMemAddr {
		return &pb.BeginReceiveResponse{Initiated: false},
			status.Errorf(codes.InvalidArgument,
				"invalid argument: dst address range [0x%x, 0x%x) not in GPU memory range [0x%x, 0x%x)",
				start, end, gpu.minMemAddr, gpu.maxMemAddr)
	}

	// protobuf default value issue
	var srcRank uint32 = 0
	if req.SrcRank != nil {
		srcRank = req.SrcRank.Value
	}

	var streamId uint64 = 0
	if req.StreamId != nil {
		streamId = req.StreamId.Value
	}

	// acquire lock on monitor
	gpu.streamDst[srcRank].cond.L.Lock()
	defer gpu.streamDst[srcRank].cond.L.Unlock()

	// assign destination info
	gpu.streamDst[srcRank].info[streamId] = &StreamDstInfo{DstAddr: req.RecvBuffAddr.Value, op: req.Op}

	// update status
	gpu.streamStatus[srcRank].Store(streamId, pb.Status_IN_PROGRESS)

	// wake up waiting threads
	log.Printf("info is: %v (len %d)", gpu.streamDst[srcRank].info, len(gpu.streamDst[srcRank].info))
	log.Printf("notifying threads waiting on streamId: %d, srcRank: %d", streamId, srcRank)
	gpu.streamDst[srcRank].cond.Broadcast()

	return &pb.BeginReceiveResponse{
		Initiated: true,
	}, nil
}

func (gpu *GPUDeviceServer) StreamSend(
	stream pb.GPUDevice_StreamSendServer,
) error {
	gpu.resetLock.RLock()
	defer gpu.resetLock.RUnlock()

	shouldError := gpu.fi.MaybeInject()
	if shouldError {
		return status.Error(
			codes.Internal,
			"GPU: (injected) internal error!",
		)
	}

	var srcRank uint32 = 0
	var streamId uint64 = 0

	for {
		// Receive a DataChunk message from the stream
		req, err := stream.Recv()
		log.Printf("stream send req is %v", req)
		// protobuf default value / stream handling, it's weird
		if req != nil {
			if req.SrcRank != nil {
				srcRank = req.SrcRank.Value
			}
			if req.StreamId != nil {
				streamId = req.StreamId.Value
			}
		}

		if err == nil {
			gpu.streamDst[srcRank].cond.L.Lock()
			info, exists := gpu.streamDst[srcRank].info[streamId]
			for !exists {
				log.Printf("waiting on recv for stream with streamId: %d, srcRank: %d", streamId, srcRank)
				log.Printf("info while waiting: %v (exists: %v)", gpu.streamDst[srcRank].info, exists)
				gpu.streamDst[srcRank].cond.Wait()
				info, exists = gpu.streamDst[srcRank].info[streamId]
			}

			log.Printf("got recv for stream with streamId: %d, srcRank: %d", streamId, srcRank)
			addr := info.DstAddr
			info.DstAddr += 8
			op := gpu.streamDst[srcRank].info[streamId].op
			gpu.streamDst[srcRank].cond.L.Unlock() // technically not thread safe

			srcData, _ := gpu.memory.Read(addr, uint64(len(req.Data)))
			gpu.memory.Write(addr, reduce(op, srcData, req.Data))

		} else if err == io.EOF {
			gpu.streamStatus[srcRank].Store(streamId, pb.Status_SUCCESS)

			// End of stream, send response
			response := &pb.StreamSendResponse{
				Success: true,
			}
			log.Printf("End of stream with streamId %d, rank %d", streamId, srcRank)
			return stream.SendAndClose(response)

		} else {
			// Handle errors during streaming
			gpu.streamStatus[srcRank].Store(streamId, pb.Status_FAILED)
			return status.Errorf(codes.Internal, "failed to receive stream: %v", err)
		}
	}
}

func (gpu *GPUDeviceServer) GetStreamStatus(
	ctx context.Context,
	req *pb.GetStreamStatusRequest,
) (*pb.GetStreamStatusResponse, error) {
	gpu.resetLock.RLock()
	defer gpu.resetLock.RUnlock()
	shouldError := gpu.fi.MaybeInject()
	if shouldError {
		return nil, status.Error(
			codes.Internal,
			"GPU: (injected) internal error!",
		)
	}

	// protobuf defaults again...
	var streamId uint64 = 0
	var srcRank uint32 = 0
	if req.StreamId != nil {
		streamId = req.StreamId.Value
	}
	if req.SrcRank != nil {
		srcRank = req.SrcRank.Value
	}
	log.Printf("status request: %v", req)

	streamStatus, ok := gpu.streamStatus[srcRank].Load(streamId)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "stream ID %d not found", streamId)
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

func (gpu *GPUDeviceServer) StreamSendDispatch() {
	log.Print("started streamsend dispatch with channel ", &gpu.streamSrc)
	for streamInfo := range gpu.streamSrc {
		streamId, srcAddr, numBytes, dstRank := streamInfo.StreamId, streamInfo.SendBuffAddr, streamInfo.NumBytes, streamInfo.DstRank
		go gpu.StreamSendExecute(streamId, srcAddr, numBytes, dstRank)
	}
}

func (gpu *GPUDeviceServer) StreamSendExecute(
	streamId, srcAddr, numBytes uint64,
	dstRank uint32,
) {
	// what to do in error cases?
	rpc_client := gpu.peers[dstRank]

	stream, err := rpc_client.StreamSend(context.Background())
	if err != nil {
		log.Printf("GPUDevice failed to start stream to rank %d", dstRank)
		return
	}

	log.Printf("starting streamsendexecute with %d, %x, %x, %d", streamId, srcAddr, numBytes, dstRank)
	// not streaming in chunks. in cuda, you specify a datatype, which makes this make sense

	var offset uint64
	for offset = 0; offset < numBytes; offset += 8 {
		data, _ := gpu.memory.Read(srcAddr+offset, 8)
		dc := pb.DataChunk{
			Data:     data,
			StreamId: &pb.StreamId{Value: streamId},
			SrcRank:  &pb.Rank{Value: gpu.rank},
		}
		log.Printf("offset: %d, len data: %d", offset, len(dc.Data))
		stream.Send(&dc)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Printf("received error: %v", err)
		log.Printf("GPUDevice failed to receive stream response from rank %d", dstRank)
	}
	gpu.streamStatus[gpu.rank].Store(streamId, pb.Status_SUCCESS)
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

func (gpu *GPUDeviceServer) Memcpy(
	ctx context.Context,
	req *pb.MemcpyRequest,
) (*pb.MemcpyResponse, error) {
	gpu.resetLock.RLock()
	defer gpu.resetLock.RUnlock()

	shouldError := gpu.fi.MaybeInject()
	if shouldError {
		return nil, status.Error(
			codes.Internal,
			"GPU: (injected) internal error!",
		)
	}

	switch x := req.Either.(type) {
	case *pb.MemcpyRequest_HostToDevice:
		err := gpu.hostToDevice(x.HostToDevice)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "HostToDevice failed: %v", err)
		}
		return &pb.MemcpyResponse{
			Either: &pb.MemcpyResponse_HostToDevice{
				HostToDevice: &pb.MemcpyHostToDeviceResponse{Success: true},
			},
		}, nil

	case *pb.MemcpyRequest_DeviceToHost:
		data, err := gpu.deviceToHost(x.DeviceToHost)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "DeviceToHost failed: %v", err)
		}
		return &pb.MemcpyResponse{
			Either: &pb.MemcpyResponse_DeviceToHost{
				DeviceToHost: &pb.MemcpyDeviceToHostResponse{DstData: data},
			},
		}, nil

	default:
		return nil, status.Error(codes.InvalidArgument, "invalid Memcpy request type")
	}
}

func (gpu *GPUDeviceServer) hostToDevice(req *pb.MemcpyHostToDeviceRequest) error {
	start := req.DstMemAddr.Value
	end := start + uint64(len(req.HostSrcData))

	if start < gpu.minMemAddr || end > gpu.maxMemAddr {
		return fmt.Errorf("destination address range [0x%x, 0x%x) is out of GPU memory range [0x%x, 0x%x)",
			start, end, gpu.minMemAddr, gpu.maxMemAddr)
	}

	return gpu.memory.Write(start, req.HostSrcData)
}

func (gpu *GPUDeviceServer) deviceToHost(req *pb.MemcpyDeviceToHostRequest) ([]byte, error) {
	start := req.SrcMemAddr.Value
	end := start + req.NumBytes

	if start < gpu.minMemAddr || end > gpu.maxMemAddr {
		return nil, fmt.Errorf("source address range [0x%x, 0x%x) is out of GPU memory range [0x%x, 0x%x)",
			start, end, gpu.minMemAddr, gpu.maxMemAddr)
	}

	return gpu.memory.Read(start, req.NumBytes)
}

func (gpu *GPUDeviceServer) ResetGpu(
	ctx context.Context,
	req *pb.ResetGpuRequest,
) (*pb.ResetGpuResponse, error) {
	gpu.resetLock.Lock()
	defer gpu.resetLock.Unlock()

	// clear everything in streamsend channel, reset dispatch thread
	close(gpu.streamSrc)
	gpu.streamSrc = make(chan StreamSrcInfo, 16)
	go gpu.StreamSendDispatch()

	for i := range gpu.nRanks {
		if i == gpu.rank {
			gpu.streamDst[i] = nil
			gpu.streamStatus[i] = nil

		} else {
			// close connection to other GPU
			gpu.conns[i].Close()

			// nil out peer
			gpu.peers[i] = nil

			// wake threads waiting on streamDst and tell them to quit? -> no you're fine, coordinator should have sent recv calls
			gpu.streamDst[i] = nil

			// nil out streamStatus
			gpu.streamStatus[i] = nil
		}
	}

	// reset streamid
	gpu.streamId.Store(0)

	return &pb.ResetGpuResponse{Success: true}, nil
}

func reduce(op pb.ReduceOp, srcData, reqData []byte) []byte {
	// both are disposable tbh
	x, y := convert.BytesToFloat64(srcData), convert.BytesToFloat64(reqData)
	switch op {
	case pb.ReduceOp_NIL:
		return convert.Float64ToBytes(y)
	case pb.ReduceOp_SUM:
		return convert.Float64ToBytes(x + y)
	case pb.ReduceOp_PROD:
		return convert.Float64ToBytes(x * y)
	case pb.ReduceOp_MIN:
		if x < y {
			return convert.Float64ToBytes(x)
		} else {
			return convert.Float64ToBytes(y)
		}
	case pb.ReduceOp_MAX:
		if x > y {
			return convert.Float64ToBytes(x)
		} else {
			return convert.Float64ToBytes(y)
		}
	default:
		return convert.Float64ToBytes(0)
	}
}
