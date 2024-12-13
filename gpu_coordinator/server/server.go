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
	port       = flag.Int("port", 8082, "The server port")
	deviceList = flag.String(
		"device-list",
		"devices.json",
		"device list for GPUCoordinator to parse",
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
		GPUDeviceList: *deviceList,
	})

	if err != nil {
		log.Fatalf("failed to start server: %q", err)
	}
	pb.RegisterGPUCoordinatorServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
