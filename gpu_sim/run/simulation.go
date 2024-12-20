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

	commId := commInitResponse.CommId
	statusResponse, err := coordinatorClient.GetCommStatus(ctx, &pb.GetCommStatusRequest{
		CommId: commId,
	})
	if err != nil {
		log.Fatalf("GetCommStatus failed: %v", err)
	}

	fmt.Printf("Communicator %d status: %v\n", commId, statusResponse.Status)
}
