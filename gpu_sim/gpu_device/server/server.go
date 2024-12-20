package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
	sl "github.com/JoshuaMBa/dsml/gpu_sim/gpu_device/server_lib"
	logging "github.com/JoshuaMBa/dsml/logging"
	"google.golang.org/grpc"
)

var (
	port        = flag.Int("port", 8081, "The server port")
	seed        = flag.Int64("seed", 42, "Random seed for generating database data")
	sleepNs     = flag.Int64("sleep-ns", 0, "Injected latency on each request")
	failureRate = flag.Int64(
		"failure-rate",
		0,
		"Injected failure rate N (0 means no injection; o/w errors one in N requests",
	)
	responseOmissionRate = flag.Int64(
		"response-omission-rate",
		0,
		"Injected response omission rate N (0 means no injection; o/w errors one in N requests",
	)
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(logging.MakeMiddleware(logging.MakeLogger())))
	server, err := sl.MakeGPUDeviceServer(sl.GPUDeviceOptions{
		SleepNs:              *sleepNs,
		FailureRate:          *failureRate,
		ResponseOmissionRate: *responseOmissionRate,
	})

	if err != nil {
		log.Fatalf("failed to start server: %q", err)
	}
	pb.RegisterGPUDeviceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
