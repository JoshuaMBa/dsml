package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
	util "github.com/JoshuaMBa/dsml/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestAllReduceRingBasic(t *testing.T) {
	coordinatorAddr := "127.0.0.1:6000"

	coordinatorConn, err := grpc.Dial(coordinatorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to GPUCoordinator: %v", err)
	}
	defer coordinatorConn.Close()
	coordinatorClient := pb.NewGPUCoordinatorClient(coordinatorConn)

	numDevices := uint32(3)
	log.Printf("Requesting a communicator with %d devices...\n", numDevices)
	ctx := context.Background()

	commInitResponse, err := coordinatorClient.CommInit(ctx, &pb.CommInitRequest{
		NumDevices: numDevices,
	})
	if err != nil {
		log.Fatalf("CommInit failed: %v", err)
	}

	if commInitResponse.Success {
		log.Printf("CommInit succeeded. CommId: %d\n", commInitResponse.CommId)
	} else {
		log.Fatalf("CommInit failed.")
	}

	log.Printf("Initializing data on GPUs...\n")
	expected := make([]float64, numDevices)
	memAddrs := make(map[uint32]*pb.MemAddr)
	for i, data := range commInitResponse.Devices {
		memAddrs[uint32(i)] = &pb.MemAddr{
			Value: data.MinMemAddr.Value,
		}

		var bytes []byte
		for j := 0; j < int(numDevices); j++ {
			expected[j] += float64(i + j)
			bytes = append(bytes, util.Float64ToBytes(float64(i+j))...)
		}
		_, err := coordinatorClient.Memcpy(ctx, &pb.MemcpyRequest{
			Either: &pb.MemcpyRequest_HostToDevice{
				HostToDevice: &pb.MemcpyHostToDeviceRequest{
					HostSrcData: bytes,
					DstDeviceId: &pb.DeviceId{Value: data.DeviceId.Value},
					DstMemAddr:  &pb.MemAddr{Value: data.MinMemAddr.Value},
				},
			},
		})

		if err != nil {
			log.Fatalf("Failed to copy initial data to device of rank 0")
		}
	}

	log.Printf("Initiating Ring Reduce...\n")
	_, err = coordinatorClient.AllReduceRing(
		ctx,
		&pb.AllReduceRingRequest{
			CommId:   commInitResponse.CommId,
			Count:    uint64(numDevices),
			Op:       pb.ReduceOp_SUM,
			MemAddrs: memAddrs,
		},
	)
	if err != nil {
		log.Fatalf("AllReduceRing failed")
	}

	wordSize := uint64(8)
	res, err := coordinatorClient.Memcpy(ctx, &pb.MemcpyRequest{
		Either: &pb.MemcpyRequest_DeviceToHost{
			DeviceToHost: &pb.MemcpyDeviceToHostRequest{
				NumBytes:    uint64(numDevices) * wordSize,
				SrcDeviceId: &pb.DeviceId{Value: commInitResponse.Devices[0].DeviceId.Value},
				SrcMemAddr:  &pb.MemAddr{Value: commInitResponse.Devices[0].MinMemAddr.Value},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to retrieve final data from device of rank 0 with deviceId %v", commInitResponse.Devices[0].DeviceId.Value)
	}

	s := "got:"
	for i := 0; i < int(numDevices); i++ {
		f := res.GetDeviceToHost().DstData[uint64(i)*wordSize : (uint64(i)*wordSize)+wordSize]
		s += fmt.Sprintf("%v ", util.BytesToFloat64(f))
	}
	log.Printf(s + "\n")

	s = "expected:"
	for i := 0; i < int(numDevices); i++ {
		s += fmt.Sprintf("%v ", expected[i])
	}
	log.Printf(s + "\n")
}
