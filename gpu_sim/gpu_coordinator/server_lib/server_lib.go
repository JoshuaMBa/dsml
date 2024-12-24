package server_lib

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"log"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	// "google.golang.org/grpc/status"
)

type GPUCoordinatorOptions struct {
	GPUDeviceList string
}

type Operation struct {
	Execute func()
}

type Communicator struct {
	nGpus       uint64
	connections []*grpc.ClientConn
	using       []uint32
	status      pb.Status
	grouped     bool
	group       []Operation
}

type GPUCoordinatorServer struct {
	pb.UnimplementedGPUCoordinatorServer
	options GPUCoordinatorOptions

	////////////////////////////////
	// bonuses
	////////////////////////////////
	nGpus      uint32
	gpuInfos   *GPUDeviceList
	available  []uint32
	comms      map[uint64]Communicator
	nextCommId uint64
	mu         sync.Mutex // General mutex
}

func MakeGPUCoordinatorServer(options GPUCoordinatorOptions) (*GPUCoordinatorServer, error) {
	gpuInfos, _ := ParseJSONFile(options.GPUDeviceList)
	server := &GPUCoordinatorServer{
		options:    options,
		nGpus:      uint32(len(gpuInfos.GPUDevices)),
		gpuInfos:   gpuInfos,
		comms:      make(map[uint64]Communicator),
		nextCommId: 0,
	}

	for i := range server.nGpus {
		server.available = append(server.available, i)
	}

	return server, nil
}

func makeConnectionClient(ServiceAddr string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(ServiceAddr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (server *GPUCoordinatorServer) commDestroyInternal(commId uint64) { // Not thread-safe on its own
	if comm, exists := server.comms[commId]; exists {
		for _, conn := range comm.connections {
			conn.Close()
		}
	}

	server.available = append(server.available, server.comms[commId].using...)
	delete(server.comms, commId)
}

func (server *GPUCoordinatorServer) CommDestroy(
	ctx context.Context,
	req *pb.CommDestroyRequest,
) (*pb.CommDestroyResponse, error) {
	server.mu.Lock()
	defer server.mu.Unlock()

	_, exists := server.comms[req.CommId]
	if !exists {
		return &pb.CommDestroyResponse{Success: false}, status.Error(codes.NotFound, "communicator not found")
	}

	server.commDestroyInternal(req.CommId)
	return &pb.CommDestroyResponse{Success: true}, nil
}

func (server *GPUCoordinatorServer) CommInit(
	ctx context.Context,
	req *pb.CommInitRequest,
) (*pb.CommInitResponse, error) {
	/*
	   1. check req.numdevices is valid
	   2. request metadata + make communicator for each gpudevice service
	*/
	server.mu.Lock()
	l := uint32(len(server.available))
	if req.NumDevices > l {
		return &pb.CommInitResponse{Success: false},
			status.Error(codes.OutOfRange, "GPUCoordinator: CommInit: num devices in communicator exceeded num gpus available")
	}

	var metadata []*pb.DeviceMetadata
	var devs []GPUDeviceInfo
	var using []uint32

	// Track GPUs being used by new communicator and available GPUs for new communicators
	using = append(using, server.available[l-req.NumDevices:]...)
	server.available = server.available[:l-req.NumDevices]
	for _, u := range using {
		devs = append(devs, server.gpuInfos.GPUDevices[u])
	}

	commId := server.nextCommId
	server.nextCommId++
	server.mu.Unlock()

	var connections []*grpc.ClientConn
	rankToAddress := make(map[uint32]string)

	// Create device clients for each device in the communicator
	for i, device := range devs {
		address := device.IP + ":" + strconv.Itoa(device.Port)
		rankToAddress[uint32(i)] = address

		conn, err := makeConnectionClient(address)
		if err != nil {
			server.mu.Lock()
			server.comms[commId] = Communicator{
				using:  using,
				status: pb.Status_FAILED,
			}
			server.commDestroyInternal(commId)
			server.mu.Unlock()
			log.Printf("GPUCoordinator: CommInit: failed to create connection client")
			return &pb.CommInitResponse{
				Success: false,
				CommId:  0,
			}, err
		}

		gpu := pb.NewGPUDeviceClient(conn)
		res, err := gpu.GetDeviceMetadata(ctx, &pb.GetDeviceMetadataRequest{})
		if err != nil {
			log.Printf("GPUCoordinatorServer: MakeGPUCoordinatorServer: failed to retrieve metadata on device (addr: %v) ", address)
			server.mu.Lock()
			server.comms[commId] = Communicator{
				using:  using,
				status: pb.Status_FAILED,
			}
			server.commDestroyInternal(commId)
			server.mu.Unlock()
			return &pb.CommInitResponse{
				Success: false,
				CommId:  0,
			}, err
		}
		metadata = append(metadata, res.Metadata)
		connections = append(connections, conn)
	}

	for i, conn := range connections {
		gpu := pb.NewGPUDeviceClient(conn)
		rank := uint32(i)
		_, err := gpu.SetupCommunication(ctx, &pb.SetupCommunicationRequest{
			RankToAddress: rankToAddress,
			Rank:          &pb.Rank{Value: rank},
		})
		if err != nil {
			server.mu.Lock()
			server.comms[commId] = Communicator{
				using:  using,
				status: pb.Status_FAILED,
			}
			server.commDestroyInternal(commId)
			server.mu.Unlock()
			log.Printf("GPUCoordinator: CommInit: failed to set up communication with device of rank %v", rank)
			return &pb.CommInitResponse{
				Success: false,
				CommId:  0,
			}, err
		}
	}

	server.mu.Lock()
	server.comms[commId] = Communicator{
		nGpus:       uint64(req.NumDevices),
		connections: connections,
		using:       using,
		status:      pb.Status_SUCCESS,
		grouped:     false,
	}
	server.mu.Unlock()

	res := &pb.CommInitResponse{
		Success: true,
		CommId:  commId,
		Devices: metadata,
	}

	return res, nil
}

func (server *GPUCoordinatorServer) GetCommStatus(
	ctx context.Context,
	req *pb.GetCommStatusRequest,
) (*pb.GetCommStatusResponse, error) {
	server.mu.Lock()
	defer server.mu.Unlock()
	return &pb.GetCommStatusResponse{
		Status: server.comms[req.CommId].status,
	}, nil
}

func (server *GPUCoordinatorServer) GroupStart(
	ctx context.Context,
	req *pb.GroupStartRequest,
) (*pb.GroupStartResponse, error) {
	server.mu.Lock()
	defer server.mu.Unlock()
	comm := server.comms[req.CommId]
	comm.grouped = true
	server.comms[req.CommId] = comm
	return &pb.GroupStartResponse{
		Success: true,
	}, nil
}

func (server *GPUCoordinatorServer) GroupEnd(
	ctx context.Context,
	req *pb.GroupEndRequest,
) (*pb.GroupEndResponse, error) {
	server.mu.Lock()
	comm := server.comms[req.CommId]
	comm.grouped = false
	server.comms[req.CommId] = comm
	server.mu.Unlock()

	// Now we need to iterate through comm.group and execute each operation

	return &pb.GroupEndResponse{
		Success: true,
	}, nil
}

func (server *GPUCoordinatorServer) waitForStream(
	ctx context.Context,
	streamId uint64,
	srcRank uint32,
	me pb.GPUDeviceClient,
) error {
	for {
		select {
		case <-ctx.Done():
			return status.Errorf(codes.Unknown, "waitForStream: context termination: streamId: %v", streamId)
		default:
			// Check status of current stream
			response, err := me.GetStreamStatus(
				ctx,
				&pb.GetStreamStatusRequest{
					StreamId: &pb.StreamId{
						Value: streamId,
					},
					SrcRank: &pb.Rank{
						Value: srcRank,
					},
				},
			)

			if response.Status == pb.Status_SUCCESS {
				return nil
			} else if response.Status == pb.Status_FAILED {
				log.Printf("GPUCoordinator: waitForStream: failed to get stream status of streamId %v", streamId)
				return err
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func (server *GPUCoordinatorServer) rebootCommunicator(commId uint64) {
	// TODO: Implement restarting the communicator (i.e. close all current streams)
}

func (server *GPUCoordinatorServer) AllReduceRing(
	ctx context.Context,
	req *pb.AllReduceRingRequest,
) (*pb.AllReduceRingResponse, error) {
	// Grouped requests are queued to be executed later
	server.mu.Lock()
	comm := server.comms[req.CommId]
	if comm.grouped {
		comm.group = append(comm.group, Operation{
			Execute: func() { server.AllReduceRing(ctx, req) },
		})
		server.comms[req.CommId] = comm
		server.mu.Unlock()
		return &pb.AllReduceRingResponse{
			Success: true,
		}, nil
	}
	server.mu.Unlock()

	// Channels through which stream ids for sent data are communicated
	recvStreamIds := make([]chan uint64, comm.nGpus)
	for i := uint64(0); i < comm.nGpus; i++ {
		recvStreamIds[i] = make(chan uint64, comm.nGpus)
	}

	// Metadata on blocks being passed around ring
	blocks := comm.nGpus
	blockSize := uint64(req.Count / comm.nGpus)
	wordSize := uint64(8)

	var wg sync.WaitGroup
	var failure atomic.Bool
	for rank := uint32(0); rank < uint32(comm.nGpus); rank++ {
		wg.Add(1)

		go func(rank uint32) {
			defer wg.Done()

			log.Printf("GPU of rank %v entered ring\n", rank)

			var prev, next uint32
			var sendBlockIndex, recvBlockIndex uint64
			if rank == 0 {
				prev = uint32(comm.nGpus - 1)
				recvBlockIndex = comm.nGpus - 1
			} else {
				prev = uint32(rank - 1)
				recvBlockIndex = uint64(rank - 1)
			}
			next = (rank + 1) % uint32(comm.nGpus)
			sendBlockIndex = uint64(rank)
			me := pb.NewGPUDeviceClient(comm.connections[rank])

			// Share-reduce phase
			for i := uint32(0); i < uint32(blocks); i++ {
				// Identify current subvector to be sent to `next`
				sendBuffAddr := req.MemAddrs[rank].Value + (sendBlockIndex * blockSize)
				numBytes := wordSize * blockSize

				// In event that comm.nGpus does not divide req.Count, add extra values to
				// last block
				if sendBlockIndex == (blocks - 1) {
					numBytes += (wordSize * (req.Count % comm.nGpus))
				}

				log.Printf("gpu %v sending %v bytes of data from addr %v + %v to gpu %v\n", rank, numBytes, sendBuffAddr, (sendBlockIndex * blockSize), next)

				// Perform non-blocking send of data to next gpu
				sendRes, _ := me.BeginSend(
					ctx,
					&pb.BeginSendRequest{
						SendBuffAddr: &pb.MemAddr{
							Value: sendBuffAddr,
						},
						NumBytes: numBytes,
						DstRank: &pb.Rank{
							Value: uint32(next),
						},
					},
				)
				if sendRes.Initiated == false {
					failure.Store(true)
					return
				}

				// Send stream id to gpu of rank `next` to initiate communication
				recvStreamIds[next] <- sendRes.StreamId.Value

				// Get stream id to receive from gpu of rank `prev`
				recvStreamId := <-recvStreamIds[rank]

				// Identify next subvector to be received from `prev`
				recvBuffAddr := req.MemAddrs[rank].Value + (recvBlockIndex * blockSize)
				numBytes = wordSize * blockSize
				if recvBlockIndex == (blocks - 1) {
					numBytes += (wordSize * (req.Count % comm.nGpus))
				}

				// On last iteration, no reduction happens
				var op pb.ReduceOp
				if i == uint32(blocks-1) {
					op = pb.ReduceOp_NIL
				} else {
					op = req.Op
				}

				log.Printf("gpu %v receiving %v bytes of data into addr %v + %v from gpu %v\n", rank, numBytes, recvBuffAddr, (recvBlockIndex * blockSize), prev)

				// Perform non-blocking receive of data from previous GPU
				recvRes, _ := me.BeginReceive(
					ctx,
					&pb.BeginReceiveRequest{
						StreamId: &pb.StreamId{
							Value: recvStreamId,
						},
						RecvBuffAddr: &pb.MemAddr{
							Value: recvBuffAddr,
						},
						NumBytes: numBytes,
						SrcRank: &pb.Rank{
							Value: uint32(prev),
						},
						Op: op,
					},
				)
				if recvRes.Initiated == false {
					failure.Store(true)
					return
				}

				// Wait for current iteration of ring reduction to complete
				srcRank := prev
				if err := server.waitForStream(ctx, sendRes.StreamId.Value, srcRank, me); err != nil {
					log.Printf("GPUCoordinator: waitForStream: failed receive (commId: %v, srcRank: %v, dstRank: %v, streamId: %v)", req.CommId, prev, rank, recvStreamId)
					failure.Store(true)
					return
				}

				srcRank = rank
				if err := server.waitForStream(ctx, recvStreamId, srcRank, me); err != nil {
					log.Printf("GPUCoordinator: waitForStream: failed send (commId: %v, srcRank: %v, dstRank: %v, streamId: %v)", req.CommId, rank, next, recvStreamId)
					failure.Store(true)
					return
				}

				// Send most recently received buffer in next iteration
				sendBlockIndex = recvBlockIndex
				if recvBlockIndex == 0 {
					recvBlockIndex = comm.nGpus - 1
				} else {
					recvBlockIndex--
				}
			}

			// Share-only phase
			for i := uint32(0); i < uint32(blocks); i++ {
				// Identify current subvector to be sent to `next`
				sendBuffAddr := req.MemAddrs[rank].Value + (sendBlockIndex * blockSize)
				numBytes := wordSize * blockSize

				// In event that comm.nGpus does not divide req.Count, add extra values to
				// last block
				if sendBlockIndex == (blocks - 1) {
					numBytes += (wordSize * (req.Count % comm.nGpus))
				}

				log.Printf("gpu %v sending %v bytes of data from addr %v + %v to gpu %v\n", rank, numBytes, sendBuffAddr, (sendBlockIndex * blockSize), next)

				// Perform non-blocking send of data to next gpu
				sendRes, _ := me.BeginSend(
					ctx,
					&pb.BeginSendRequest{
						SendBuffAddr: &pb.MemAddr{
							Value: sendBuffAddr,
						},
						NumBytes: numBytes,
						DstRank: &pb.Rank{
							Value: uint32(next),
						},
					},
				)
				if sendRes.Initiated == false {
					failure.Store(true)
					return
				}

				// Send stream id to gpu of rank `next` to initiate communication
				recvStreamIds[next] <- sendRes.StreamId.Value

				// Get stream id to receive from gpu of rank `prev`
				recvStreamId := <-recvStreamIds[rank]

				// Identify next subvector to be received from `prev`
				recvBuffAddr := req.MemAddrs[rank].Value + (recvBlockIndex * blockSize)
				numBytes = wordSize * blockSize
				if recvBlockIndex == (blocks - 1) {
					numBytes += (wordSize * (req.Count % comm.nGpus))
				}

				log.Printf("gpu %v receiving %v bytes of data into addr %v + %v from gpu %v\n", rank, numBytes, recvBuffAddr, (recvBlockIndex * blockSize), prev)

				// Perform non-blocking receive of data from previous GPU
				recvRes, _ := me.BeginReceive(
					ctx,
					&pb.BeginReceiveRequest{
						StreamId: &pb.StreamId{
							Value: recvStreamId,
						},
						RecvBuffAddr: &pb.MemAddr{
							Value: recvBuffAddr,
						},
						NumBytes: numBytes,
						SrcRank: &pb.Rank{
							Value: uint32(prev),
						},
						Op: pb.ReduceOp_NIL,
					},
				)
				if recvRes.Initiated == false {
					failure.Store(true)
					return
				}

				// Wait for current iteration of ring reduction to complete
				srcRank := prev
				if err := server.waitForStream(ctx, sendRes.StreamId.Value, srcRank, me); err != nil {
					log.Printf("GPUCoordinator: waitForStream: failed receive (commId: %v, srcRank: %v, dstRank: %v, streamId: %v)", req.CommId, prev, rank, recvStreamId)
					failure.Store(true)
					return
				}

				srcRank = rank
				if err := server.waitForStream(ctx, recvStreamId, srcRank, me); err != nil {
					log.Printf("GPUCoordinator: waitForStream: failed send (commId: %v, srcRank: %v, dstRank: %v, streamId: %v)", req.CommId, rank, next, recvStreamId)
					failure.Store(true)
					return
				}

				// Send most recently received buffer in next iteration
				sendBlockIndex = recvBlockIndex
				if recvBlockIndex == 0 {
					recvBlockIndex = comm.nGpus - 1
				} else {
					recvBlockIndex--
				}
			}
		}(rank)
	}

	wg.Wait()
	if failure.Load() == true {
		server.rebootCommunicator(req.CommId)
		return &pb.AllReduceRingResponse{
			Success: false,
		}, status.Errorf(codes.Unknown, "AllReduceRing: ring crash: commId: %v", req.CommId)
	}

	return &pb.AllReduceRingResponse{
		Success: true,
	}, nil
}

func (server *GPUCoordinatorServer) Memcpy(
	ctx context.Context,
	req *pb.MemcpyRequest,
) (*pb.MemcpyResponse, error) {
	server.mu.Lock()
	defer server.mu.Unlock()

	var targetDeviceId uint64
	switch x := req.Either.(type) {
	case *pb.MemcpyRequest_HostToDevice:
		targetDeviceId = x.HostToDevice.DstDeviceId.Value
	case *pb.MemcpyRequest_DeviceToHost:
		targetDeviceId = x.DeviceToHost.SrcDeviceId.Value
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid Memcpy request type")
	}

	fmt.Printf("Memcpy called for targetDeviceId: %d\n", targetDeviceId)

	var targetGPU pb.GPUDeviceClient
	for cid, comm := range server.comms {
		fmt.Printf("Comm ID: %d\n", cid)
		for i, u := range comm.using {
			deviceID := uint64(server.gpuInfos.GPUDevices[u].DeviceID)
			fmt.Printf("Checking DeviceID: %d against targetDeviceId: %d\n", deviceID, targetDeviceId)
			if uint64(server.gpuInfos.GPUDevices[u].DeviceID) == targetDeviceId {
				targetGPU = pb.NewGPUDeviceClient(comm.connections[i])
				break
			}
		}
		if targetGPU != nil {
			break
		}
	}
	fmt.Printf("Memcpy called for targetDeviceId: %d\n", targetDeviceId)

	if targetGPU == nil {
		return nil, status.Error(codes.NotFound, "target device not found")
	}

	// Forward the Memcpy request to the target GPU device
	resp, err := targetGPU.Memcpy(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Memcpy operation failed: %v", err)
	}

	return resp, nil
}

////////////////////////////////////////////
// helpers \ bookkeeping
////////////////////////////////////////////

type GPUDeviceInfo struct {
	DeviceID int    `json:"device_id"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
}

type GPUDeviceList struct {
	GPUDevices []GPUDeviceInfo `json:"gpu_devices"`
}

func ParseJSONFile(filePath string) (*GPUDeviceList, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config GPUDeviceList
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
