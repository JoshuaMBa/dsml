package server_lib

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	// "log"
	// "sort"
	"sync"
	// "sync/atomic"
	// "time"

	"github.com/JoshuaMBa/dsml/gpu_coordinator/proto"
	cpb "github.com/JoshuaMBa/dsml/gpu_coordinator/proto"
	dpb "github.com/JoshuaMBa/dsml/gpu_device/proto"

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
	nGpus   uint64
	gpus    []*dpb.GPUDeviceClient
	using   []uint32
	status  cpb.CommStatus
	grouped bool
	group   []Operation
}

type GPUCoordinatorServer struct {
	proto.UnimplementedGPUCoordinatorServer
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
		nextCommId: 0,
	}

	for i := range server.nGpus {
		server.available = append(server.available, i)
	}

	return server, nil
}

func makeConnectionClient(ServiceAddr string) (*dpb.GPUDeviceClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(ServiceAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := dpb.NewGPUDeviceClient(conn)

	return &client, nil
}

func (server *GPUCoordinatorServer) CommInit(
	ctx context.Context,
	req *cpb.CommInitRequest,
) (*cpb.CommInitResponse, error) {
	/*
	   1. check req.numdevices is valid
	   2. request metadata + make communicator for each gpudevice service
	*/
	server.mu.Lock()
	l := uint32(len(server.available))
	if req.NumDevices > l {
		return &cpb.CommInitResponse{Success: false},
			status.Error(codes.OutOfRange, "num devices in communicator exceeded num gpus available")
	}

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

	var gpus []*dpb.GPUDeviceClient

	// Create device clients for each device in the communicator
	for i := range req.NumDevices {
		device := devs[i]

		gpu, err := makeConnectionClient(device.IP + ":" + strconv.Itoa(device.Port))
		if err != nil {
			server.mu.Lock()
			server.available = append(server.available, using...)
			server.mu.Unlock()
			return &cpb.CommInitResponse{
				Success: false,
				CommId:  0,
			}, err
		}

		gpus = append(gpus, gpu)
	}

	res := &cpb.CommInitResponse{
		Success: true,
		CommId:  commId,
	}

	server.mu.Lock()
	server.comms[commId] = Communicator{
		nGpus:   uint64(req.NumDevices),
		gpus:    gpus,
		using:   using,
		status:  cpb.CommStatus_SUCCESS,
		grouped: false,
	}
	server.mu.Unlock()

	return res, nil
}

func (server *GPUCoordinatorServer) GetCommStatus(
	ctx context.Context,
	req *cpb.GetCommStatusRequest,
) (*cpb.GetCommStatusResponse, error) {
	server.mu.Lock()
	defer server.mu.Unlock()
	return &cpb.GetCommStatusResponse{
		Status: server.comms[req.CommId].status,
	}, nil
}

func (server *GPUCoordinatorServer) GroupStart(
	ctx context.Context,
	req *cpb.GroupStartRequest,
) (*cpb.GroupStartResponse, error) {
	server.mu.Lock()
	defer server.mu.Unlock()
	comm := server.comms[req.CommId]
	comm.grouped = true
	server.comms[req.CommId] = comm
	return &cpb.GroupStartResponse{
		Success: true,
	}, nil
}

func (server *GPUCoordinatorServer) GroupEnd(
	ctx context.Context,
	req *cpb.GroupEndRequest,
) (*cpb.GroupEndResponse, error) {
	server.mu.Lock()
	comm := server.comms[req.CommId]
	comm.grouped = false
	server.comms[req.CommId] = comm
	server.mu.Unlock()

	// Now we need to iterate through comm.group and execute each operation

	return &cpb.GroupEndResponse{
		Success: true,
	}, nil
}

func (server *GPUCoordinatorServer) AllReduceRing(
	ctx context.Context,
	req *cpb.AllReduceRingRequest,
) (*cpb.AllReduceRingResponse, error) {
	// Grouped requests are queued to be executed later
	server.mu.Lock()
	comm := server.comms[req.CommId]
	if comm.grouped {
		comm.group = append(comm.group, Operation{
			Execute: func() { server.AllReduceRing(ctx, req) },
		})
		server.comms[req.CommId] = comm
		server.mu.Unlock()
		return &cpb.AllReduceRingResponse{
			Success: true,
		}, nil
	}
	server.mu.Unlock()
	panic("unimplemented")
}

func (server *GPUCoordinatorServer) Memcpy(
	ctx context.Context,
	req *cpb.MemcpyRequest,
) (*cpb.MemcpyResponse, error) {
	panic("unimplemented")
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
