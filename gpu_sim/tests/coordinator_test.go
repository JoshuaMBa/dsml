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

func AllReduceRingBasic(t *testing.T, coordinatorAddr string, numDevices uint32, op pb.ReduceOp, operation func(float64, float64) float64) {
	coordinatorConn, err := grpc.Dial(coordinatorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to GPUCoordinator: %v", err)
	}
	defer coordinatorConn.Close()
	coordinatorClient := pb.NewGPUCoordinatorClient(coordinatorConn)

	ctx := context.Background()
	commInitResponse, err := coordinatorClient.CommInit(ctx, &pb.CommInitRequest{
		NumDevices: numDevices,
	})

	if err != nil || commInitResponse.Success == false {
		t.Fatalf("Failed to initialize communicator: %v", err)
	}

	expected := make([]float64, numDevices)
	memAddrs := make(map[uint32]*pb.MemAddr)
	for i, data := range commInitResponse.Devices {
		memAddrs[uint32(i)] = &pb.MemAddr{
			Value: data.MinMemAddr.Value,
		}

		var bytes []byte
		for j := 0; j < int(numDevices); j++ {
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
			Count:    uint64(numDevices),
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
				NumBytes:    uint64(numDevices) * wordSize,
				SrcDeviceId: &pb.DeviceId{Value: commInitResponse.Devices[0].DeviceId.Value},
				SrcMemAddr:  &pb.MemAddr{Value: commInitResponse.Devices[0].MinMemAddr.Value},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to retrieve final data from device of rank 0")
	}

	var got []float64
	for i := 0; i < int(numDevices); i++ {
		f := res.GetDeviceToHost().DstData[uint64(i)*wordSize : (uint64(i)*wordSize)+wordSize]
		got = append(got, util.BytesToFloat64(f))
	}

	const tolerance = 1e-9
	for i := 0; i < int(numDevices); i++ {
		if math.Abs(expected[i]-got[i]) > tolerance {
			s := "got:"
			for j := 0; j < int(numDevices); j++ {
				s += fmt.Sprintf("%v ", got[i])
			}
			s += "\nexpected:"
			for i := 0; i < int(numDevices); i++ {
				s += fmt.Sprintf("%v ", expected[i])
			}
			t.Fatalf(s + "\n")
		}
	}
}

func TestAllReduceRingBasic(t *testing.T) {
	AllReduceRingBasic(t, "127.0.0.1:6000", 3, pb.ReduceOp_SUM, func(a float64, b float64) float64 {
		return a + b
	})
}
