package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"testing"

	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
	util "github.com/JoshuaMBa/dsml/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AllReduceRingNoFailureInjection(
	t *testing.T,
	coordinatorClient pb.GPUCoordinatorClient,
	numDevices uint32,
	dataLength uint32,
	op pb.ReduceOp,
	operation func(float64, float64) float64,
) {
	ctx := context.Background()
	commInitResponse, err := coordinatorClient.CommInit(ctx, &pb.CommInitRequest{
		NumDevices: numDevices,
	})
	if err != nil || commInitResponse.Success == false {
		t.Fatalf("Failed to initialize communicator: %v", err)
	}

	defer coordinatorClient.CommDestroy(ctx, &pb.CommDestroyRequest{CommId: commInitResponse.CommId})

	expected := make([]float64, dataLength)
	memAddrs := make(map[uint32]*pb.MemAddr)
	for i, data := range commInitResponse.Devices {
		memAddrs[uint32(i)] = &pb.MemAddr{
			Value: data.MinMemAddr.Value,
		}

		var bytes []byte
		for j := 0; j < int(dataLength); j++ {
			if i == 0 {
				expected[j] = float64(i + j)
			} else {
				expected[j] = operation(expected[j], float64(i+j))
			}
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

	allReduceRingResponse, err := coordinatorClient.AllReduceRing(
		ctx,
		&pb.AllReduceRingRequest{
			CommId:   commInitResponse.CommId,
			Count:    uint64(dataLength),
			Op:       op,
			MemAddrs: memAddrs,
		},
	)
	if err != nil || allReduceRingResponse.Success == false {
		t.Fatalf("AllReduceRing failure detected")
	}

	wordSize := uint64(8)
	res, err := coordinatorClient.Memcpy(ctx, &pb.MemcpyRequest{
		Either: &pb.MemcpyRequest_DeviceToHost{
			DeviceToHost: &pb.MemcpyDeviceToHostRequest{
				NumBytes:    uint64(dataLength) * wordSize,
				SrcDeviceId: &pb.DeviceId{Value: commInitResponse.Devices[0].DeviceId.Value},
				SrcMemAddr:  &pb.MemAddr{Value: commInitResponse.Devices[0].MinMemAddr.Value},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to retrieve final data from device of rank 0")
	}

	var got []float64
	for i := 0; i < int(dataLength); i++ {
		f := res.GetDeviceToHost().DstData[uint64(i)*wordSize : (uint64(i)*wordSize)+wordSize]
		got = append(got, util.BytesToFloat64(f))
	}

	const tolerance = 1e-9
	for i := 0; i < int(dataLength); i++ {
		if math.Abs(expected[i]-got[i]) > tolerance {
			s := "got:"
			for j := 0; j < int(dataLength); j++ {
				s += fmt.Sprintf("%v ", got[i])
			}
			s += "\nexpected:"
			for i := 0; i < int(dataLength); i++ {
				s += fmt.Sprintf("%v ", expected[i])
			}
			t.Fatalf(s + "\n")
		}
	}
}

func connectToCoordinator(t *testing.T) *grpc.ClientConn {
	coordinatorAddr := "127.0.0.1:6000"
	coordinatorConn, err := grpc.Dial(coordinatorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to GPUCoordinator: %v", err)
	}
	return coordinatorConn
}

func AllOperationsNoFailureInjection(
	t *testing.T,
	numDevices uint32,
	dataLength uint32,
) {
	coordinatorConn := connectToCoordinator(t)
	defer coordinatorConn.Close()
	coordinatorClient := pb.NewGPUCoordinatorClient(coordinatorConn)

	AllReduceRingNoFailureInjection(t, coordinatorClient, numDevices, dataLength, pb.ReduceOp_SUM, func(a float64, b float64) float64 {
		return a + b
	})
	AllReduceRingNoFailureInjection(t, coordinatorClient, numDevices, dataLength, pb.ReduceOp_PROD, func(a float64, b float64) float64 {
		return a * b
	})
	AllReduceRingNoFailureInjection(t, coordinatorClient, numDevices, dataLength, pb.ReduceOp_MAX, func(a float64, b float64) float64 {
		if a < b {
			return b
		} else {
			return a
		}
	})
	AllReduceRingNoFailureInjection(t, coordinatorClient, numDevices, dataLength, pb.ReduceOp_MIN, func(a float64, b float64) float64 {
		if a < b {
			return a
		} else {
			return b
		}
	})
}

func TestAllReduceRingBasic(t *testing.T) {
	// GPUs pass one float around the ring
	numDevices := uint32(3)
	dataLength := uint32(3)

	AllOperationsNoFailureInjection(t, numDevices, dataLength)
}

func TestAllReduceRingCoprime(t *testing.T) {
	// Number of devices does not divide the dataLength
	numDevices := uint32(3)
	dataLength := uint32(17)

	AllOperationsNoFailureInjection(t, numDevices, dataLength)
}
