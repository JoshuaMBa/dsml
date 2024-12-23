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

func loadConfigFromJSON(filePath string) (int, *sl.GPUDeviceOptions, error) {
	// Configuration struct for JSON parsing
	type Config struct {
		Port             int                 `json:"port"`
		GPUDeviceOptions sl.GPUDeviceOptions `json:"gpu_device_options"`
	}

	file, err := os.Open(filePath)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return 0, nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return config.Port, &config.GPUDeviceOptions, nil
}

func main() {
	flag.Parse()

	var options *sl.GPUDeviceOptions
	var serverPort int

	if *configFile != "" {
		// Load configuration from JSON
		var err error
		serverPort, options, err = loadConfigFromJSON(*configFile)
		if err != nil {
			log.Fatalf("failed to load configuration: %v", err)
		}
	} else {
		// Use command-line flags if JSON file is not provided
		serverPort = *port
		options = sl.DefaultGPUDeviceOptions()
		options.SleepNs = *sleepNs
		options.FailureRate = *failureRate
		options.ResponseOmissionRate = *responseOmissionRate
	}

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create the server
	s := grpc.NewServer(grpc.UnaryInterceptor(logging.MakeMiddleware(logging.MakeLogger())))
	server, err := sl.MakeGPUDeviceServer(*options)
	if err != nil {
		log.Fatalf("failed to start server: %q", err)
	}

	go server.StreamSendThread()

	pb.RegisterGPUDeviceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
