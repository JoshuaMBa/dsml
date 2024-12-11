package server_lib

import (
	"context"
	// "log"
	// "sort"
	// "sync"
	// "sync/atomic"
	// "time"
    
    cpb "github.com/JoshuaMBa/dsml/gpu_coordinator/proto"
	// dpb "github.com/JoshuaMBa/dsml/gpu_device/proto"

    // "google.golang.org/grpc"
    // "google.golang.org/grpc/codes"
    // "google.golang.org/grpc/credentials"
    // "google.golang.org/grpc/status"
)

type GPUCoordinatorOptions{
    // Server address for the GPUDevice
    GPUDeviceAddr string
}

type GPUCoordinatorServer struct {
    options GPUCoordinatorOptions 
}

func MakeGPUCoordinatorServer(options GPUCoordinatorOptions) (*GPUCoordinatorServer, error) {
    server := &GPUCoordinatorServer{
        options: options,
    }

    return server, nil
}

func (server *GPUCoordinatorServer) CommInit(
    ctx context.Context, 
    req *cpb.CommInitRequest,
) (cpb.CommInitResponse, error) {
    
}

func (server *GPUCoordinatorServer) GetCommStatus(
    ctx context.Context,
    req *cpb.GetCommStatusRequest,
) (*cpb.GetCommStatusResponse, error) {

}

func (server *GPUCoordinatorServer) GroupStart(
    ctx context.Context,
    req *cpb.GroupStartRequest,
) (*cpb.GroupStartResponse, error) {

}

func (server *GPUCoordinatorServer) GroupEnd(
    ctx context.Context,
    req *cpb.GroupEndRequest,
) (*cpb.GroupEndResponse, error) {

}

func (server *GPUCoordinatorServer) AllReduceRing(
    ctx context.Context,
    req *cpb.AllReduceRingRequest,
) (*cpb.AllReduceRingResponse, error) {

}

func (server *GPUCoordinatorServer) Memcpy(
    ctx context.Context,
    req *cpb.MemcpyRequest,
) (*cpb.MemcpyResponse, error) {

}
