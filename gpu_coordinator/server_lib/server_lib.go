package server_lib

import (
	"context"
	"encoding/json"
	"os"

	// "log"
	// "sort"
	// "sync"
	// "sync/atomic"
	// "time"

	"github.com/JoshuaMBa/dsml/gpu_coordinator/proto"
	cpb "github.com/JoshuaMBa/dsml/gpu_coordinator/proto"
	dpb "github.com/JoshuaMBa/dsml/gpu_device/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/status"
)

type GPUCoordinatorOptions struct {
	GPUDeviceList string
}

type Communicator struct {
	nGpus uint64
	gpus  []*dpb.GPUDeviceClient
}

type GPUCoordinatorServer struct {
	proto.UnimplementedGPUCoordinatorServer
	options GPUCoordinatorOptions

	////////////////////////////////
	// bonuses
	////////////////////////////////
	nGpus    uint32
	gpuInfos *GPUDeviceList
}

func MakeGPUCoordinatorServer(options GPUCoordinatorOptions) (*GPUCoordinatorServer, error) {
	gpuInfos, _ := ParseJSONFile(options.GPUDeviceList)
	server := &GPUCoordinatorServer{
		options:  options,
		nGpus:    uint32(len(gpuInfos.GPUDevices)),
		gpuInfos: gpuInfos,
	}

	return server, nil
}

func (server *GPUCoordinatorServer) CommInit(
	ctx context.Context,
	req *cpb.CommInitRequest,
) (*cpb.CommInitResponse, error) {
	/*
	   1. check req.numdevices is valid
	   2. request metadata + make communicator for each gpudevice service
	*/
	if req.NumDevices > server.nGpus {
		return &cpb.CommInitResponse{Success: false},
			status.Error(codes.OutOfRange, "num devices in communicator exceeded num gpus available")
	}

	comm := Communicator{}
	// try to connect to each gpu in gpuinfos

	panic("unimplemented")
}

func (server *GPUCoordinatorServer) GetCommStatus(
	ctx context.Context,
	req *cpb.GetCommStatusRequest,
) (*cpb.GetCommStatusResponse, error) {
	panic("unimplemented")
}

func (server *GPUCoordinatorServer) GroupStart(
	ctx context.Context,
	req *cpb.GroupStartRequest,
) (*cpb.GroupStartResponse, error) {
	panic("unimplemented")
}

func (server *GPUCoordinatorServer) GroupEnd(
	ctx context.Context,
	req *cpb.GroupEndRequest,
) (*cpb.GroupEndResponse, error) {
	panic("unimplemented")
}

func (server *GPUCoordinatorServer) AllReduceRing(
	ctx context.Context,
	req *cpb.AllReduceRingRequest,
) (*cpb.AllReduceRingResponse, error) {
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
