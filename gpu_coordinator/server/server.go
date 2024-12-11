package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/JoshuaMBa/dsml/gpu_coordinator/proto"
	sl "github.com/JoshuaMBa/dsml/gpu_coordinator/server_lib"
	logging "github.com/JoshuaMBa/dsml/logging"
	"google.golang.org/grpc"
)

var (
	port          = flag.Int("port", 8082, "The server port")
	gpuDeviceAddr = flag.String(
		"user-service",
		"[::1]:8080",
		"Server address for the GPUDevice",
	)
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(logging.MakeMiddleware(logging.MakeLogger())))
	server, err := sl.MakeGPUCoordinatorServer(sl.GPUCoordinatorOptions{
		GPUDeviceAddr: *gpuDeviceAddr,
	})

	if err != nil {
		log.Fatalf("failed to start server: %q", err)
	}
	go server.ContinuallyRefreshCache()
	pb.RegisterGPUCoordinatorServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	server.Close()
}
