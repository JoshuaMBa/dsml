package server_lib

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/JoshuaMBa/dsml/failure_injection"
	fipb "github.com/JoshuaMBa/dsml/failure_injection/proto"
	"github.com/JoshuaMBa/dsml/gpu_device/proto"
	pb "github.com/JoshuaMBa/dsml/gpu_device/proto"
	"google.golang.org/grpc"
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
	// device info (i don't anticipate ever using this, maybe it goes into options?)
	////////////////////////////
	deviceId   uint64
	minMemAddr uint64
	maxMemAddr uint64

	////////////////////////////
	// gpu communication logical
	////////////////////////////
	rank   uint64 // my rank in the communicator
	nRanks uint64 // total number of gpus in the communicator

	////////////////////////////
	// gpu communications, implementation detail
	////////////////////////////
	streamId    atomic.Uint64         // my streamId when sending to others
	peers       []*pb.GPUDeviceClient // rpc handles for other gpus
	streamDests []map[uint64]uint64   // where streams i am receiving should go (lookup is peer->streamId->memAddr)
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
	fiConfig := fi.GetInjectionConfig()
	log.Printf(
		"Starting GPUDevice with failure injection config: [sleepNs: %d, failureRate: %d, responseOmissionRate: %d]",
		fiConfig.SleepNs,
		fiConfig.FailureRate,
		fiConfig.ResponseOmissionRate,
	)
	return &GPUDeviceServer{
		options: options,
		fi:      fi,

		deviceId: options.DeviceId,
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
	panic("not implemented")

}

func (gpu *GPUDeviceServer) BeginReceive(
	ctx context.Context,
	req *pb.BeginReceiveRequest,
) (*pb.BeginReceiveResponse, error) {
	// get current stream id and increment
	// kick off goroutine actually sending data using streamsend
	panic("not implemented")
}

// StreamSend implements proto.GPUDeviceServer.
func (gpu *GPUDeviceServer) StreamSend(grpc.ClientStreamingServer[pb.DataChunk, pb.StreamSendResponse]) error {
	panic("unimplemented")
}

func (gpu *GPUDeviceServer) GetStreamStatus(
	ctx context.Context,
	req *pb.GetStreamStatusRequest,
) (*pb.GetStreamStatusResponse, error) {
	panic("not implemented")
}

func (db *GPUDeviceServer) SetInjectionConfig(
	ctx context.Context,
	req *fipb.SetInjectionConfigRequest,
) (*fipb.SetInjectionConfigResponse, error) {
	db.fi.SetInjectionConfigPb(req.Config)
	log.Printf(
		"GPUDevice failure injection config set to: [sleepNs: %d, failureRate: %d, responseOmissionRate: %d]",
		req.Config.SleepNs,
		req.Config.FailureRate,
		req.Config.ResponseOmissionRate,
	)

	return &fipb.SetInjectionConfigResponse{}, nil
}
