package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	coordinatorAddr := "127.0.0.1:6000"

	coordinatorConn, err := grpc.Dial(coordinatorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to GPUCoordinator: %v", err)
	}
	defer coordinatorConn.Close()
	coordinatorClient := pb.NewGPUCoordinatorClient(coordinatorConn)

	numDevices := uint32(3)
	fmt.Printf("Requesting a communicator with %d devices...\n", numDevices)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	commInitResponse, err := coordinatorClient.CommInit(ctx, &pb.CommInitRequest{
		NumDevices: numDevices,
	})
	if err != nil {
		log.Fatalf("CommInit failed: %v", err)
	}

	if commInitResponse.Success {
		fmt.Printf("CommInit succeeded. CommId: %d\n", commInitResponse.CommId)
	} else {
		fmt.Println("CommInit failed.")
		return
	}


	firstDevice := commInitResponse.Devices[0]

	commId := commInitResponse.CommId
	statusResponse, err := coordinatorClient.GetCommStatus(ctx, &pb.GetCommStatusRequest{
		CommId: commId,
	})
	if err != nil {
		log.Fatalf("GetCommStatus failed: %v", err)
	}

	fmt.Printf("Communicator %d status: %v\n", commId, statusResponse.Status)

	hostToDeviceReq := &pb.MemcpyRequest{
		Either: &pb.MemcpyRequest_HostToDevice{
			HostToDevice: &pb.MemcpyHostToDeviceRequest{
				HostSrcData: []byte("Hello, GPU!"),
				DstDeviceId: &pb.DeviceId{Value: firstDevice.DeviceId.Value},
				DstMemAddr:  &pb.MemAddr{Value: firstDevice.MinMemAddr.Value},
			},
		},
	}

	hostToDeviceResponse, err := coordinatorClient.Memcpy(ctx, hostToDeviceReq)
	if err != nil {
		log.Fatalf("Memcpy to device failed: %v", err)
	}
	log.Printf("Memcpy Host-to-Device Success: %v", hostToDeviceResponse.GetHostToDevice().Success)

	deviceToHostReq := &pb.MemcpyRequest{
		Either: &pb.MemcpyRequest_DeviceToHost{
			DeviceToHost: &pb.MemcpyDeviceToHostRequest{
				NumBytes:    uint64(12),
				SrcDeviceId: &pb.DeviceId{Value: firstDevice.DeviceId.Value},
				SrcMemAddr:  &pb.MemAddr{Value: firstDevice.MinMemAddr.Value},
			},
		},
	}
	deviceToHostResponse, err := coordinatorClient.Memcpy(ctx, deviceToHostReq)
	if err != nil {
		log.Fatalf("Memcpy from device failed: %v", err)
	}
	log.Printf("Memcpy Device-to-Host got data: %s", string(deviceToHostResponse.GetDeviceToHost().DstData))
	
	destroyResponse, err := coordinatorClient.CommDestroy(ctx, &pb.CommDestroyRequest{
		CommId: commId,
	})
	if err != nil {
		log.Fatalf("CommDestroy failed: %v", err)
	}

	if destroyResponse.Success {
		fmt.Printf("Communicator %d destroyed successfully.\n", commId)
	} else {
		fmt.Printf("Failed to destroy communicator %d.\n", commId)
	}
}
