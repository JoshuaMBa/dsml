package main

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connect2Gpu(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return grpc.NewClient(address, opts...)
}

type GPUMetadata struct {
	DeviceId   uint64
	MinMemAddr uint64
	MaxMemAddr uint64
}

func TestSendRecvBasic(t *testing.T) {
	numGpus := 3

	// connect to gpus
	gpus := make([]pb.GPUDeviceClient, numGpus)
	gpuMetadatas := make([]GPUMetadata, numGpus)

	for i := 0; i < numGpus; i++ {
		conn, _ := connect2Gpu("127.0.0.1:500" + fmt.Sprint(i+1))
		defer conn.Close()
		gpus[i] = pb.NewGPUDeviceClient(conn)

		mdr, _ := gpus[i].GetDeviceMetadata(context.Background(), &pb.GetDeviceMetadataRequest{})
		gpuMetadatas[i] = GPUMetadata{mdr.Metadata.DeviceId.Value, mdr.Metadata.MinMemAddr.Value, mdr.Metadata.MaxMemAddr.Value}

	}

	// check that gpus were setup like we expect
	for i, md := range gpuMetadatas {
		j := i + 1
		assert.Equal(t, uint64(10+j), md.DeviceId)
		assert.Equal(t, uint64(0x1000), md.MaxMemAddr-md.MinMemAddr)
	}
	fmt.Println("gpu devices setup correctly")

	// beginsend, beginrecv with 1, 2 (remember off by one, ugh)

	// check that gpus were setup like we expect
	for i, md := range gpuMetadatas {
		j := i + 1
		assert.Equal(t, uint64(10+j), md.DeviceId, "Unexpected DeviceId for GPU %d", j)
		assert.Equal(t, uint64(0x1000), md.MaxMemAddr-md.MinMemAddr, "Unexpected memory range for GPU %d", j)
	}
	t.Log("GPU devices initialized correctly")

	// group into a communicator ig
	rank2Address := make(map[uint32]string)
	for i := uint32(0); i < uint32(numGpus); i++ {
		addr := "127.0.0.1:500" + fmt.Sprint(i+1)
		rank2Address[i] = addr
	}

	for _, gpu := range gpus {
		gpu.SetupCommunication(context.Background(), &pb.SetupCommunicationRequest{RankToAddress: rank2Address})
	}

	// perform send and receive
	// resp, _ := gpus[0].BeginSend(context.Background(),
	// 	&pb.BeginSendRequest{
	// 		SendBuffAddr: &pb.MemAddr{Value: gpuMetadatas[0].MinMemAddr},
	// 		NumBytes:     0x100,
	// 		DstRank:      &pb.Rank{Value: 1},
	// 	})

	// t.Logf("streamid: %d", resp.StreamId.Value)
}
