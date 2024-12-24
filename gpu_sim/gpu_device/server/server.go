package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	sl "github.com/JoshuaMBa/dsml/gpu_sim/gpu_device/server_lib"
	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
	logging "github.com/JoshuaMBa/dsml/logging"
	"google.golang.org/grpc"
)

var (
	port        = flag.Int("port", 8081, "The server port")
	configFile  = flag.String("config", "", "Path to JSON configuration file")
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

type GPUDeviceServerConfig struct {
	Port    int
	Options *sl.GPUDeviceOptions
}

func loadConfigFromJSON(filePath string) (*GPUDeviceServerConfig, error) {
	// Configuration struct for JSON parsing
	type Config struct {
		Port             int                 `json:"port"`
		GPUDeviceOptions sl.GPUDeviceOptions `json:"gpu_device_options"`
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &GPUDeviceServerConfig{
		Port:    config.Port,
		Options: &config.GPUDeviceOptions,
	}, nil
}

func StartServer(config *GPUDeviceServerConfig) (*grpc.Server, net.Listener, error) {
	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %w", err)
	}

	// Create the server
	s := grpc.NewServer(grpc.UnaryInterceptor(logging.MakeMiddleware(logging.MakeLogger())))
	server, err := sl.MakeGPUDeviceServer(*config.Options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start server: %w", err)
	}

	pb.RegisterGPUDeviceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	return s, lis, nil
}

func main() {
	flag.Parse()

	var config *GPUDeviceServerConfig
	if *configFile != "" {
		// Load configuration from JSON
		var err error
		config, err = loadConfigFromJSON(*configFile)
		if err != nil {
			log.Fatalf("failed to load configuration: %v", err)
		}
	} else {
		// Use command-line flags if JSON file is not provided
		config = &GPUDeviceServerConfig{
			Port: *port,
			Options: &sl.GPUDeviceOptions{
				SleepNs:              *sleepNs,
				FailureRate:          *failureRate,
				ResponseOmissionRate: *responseOmissionRate,
			},
		}
	}

	server, lis, err := StartServer(config)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	defer server.GracefulStop()

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
