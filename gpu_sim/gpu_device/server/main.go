package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
	server_lib "github.com/JoshuaMBa/dsml/gpu_sim/gpu_device/server_lib"
)

func main() {
	// Configuration for the mock device
	deviceId := uint64(1)
	minMemAddr := uint64(0)
	maxMemAddr := uint64(1024 * 1024) // 1MB
	rankToAddress := map[uint32]string{
		0: "localhost:50051",
	}

	server := server_lib.MockGPUDeviceServer(deviceId, minMemAddr, maxMemAddr, rankToAddress)

	grpcServer := grpc.NewServer()
	pb.RegisterGPUDeviceServer(grpcServer, server)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("GPUDevice server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
