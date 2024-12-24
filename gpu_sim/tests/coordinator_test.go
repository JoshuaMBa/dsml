package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
	"github.com/stretchr/testify/assert"
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
	fmt.Printf("Requesting a communicator with %d devices...\n", numDevices)
	ctx := context.Background()

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

	memAddrs := make(map[uint32]*pb.MemAddr)
	for i, data := range commInitResponse.Devices {
		memAddrs[uint32(i)] = &pb.MemAddr{
			Value: data.MinMemAddr.Value,
		}
	}

	allReduceRingResponse, err := coordinatorClient.AllReduceRing(
		ctx,
		&pb.AllReduceRingRequest{
			CommId:   commInitResponse.CommId,
			Count:    3,
			Op:       pb.ReduceOp_SUM,
			MemAddrs: memAddrs,
		},
	)
	log.Printf("Error: %v", err)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, allReduceRingResponse.Success)
}
