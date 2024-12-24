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
		assert.Equal(t, uint64(10+j), md.DeviceId, "Unexpected DeviceId for GPU %d", i)
		assert.Equal(t, uint64(0x1000), md.MaxMemAddr-md.MinMemAddr, "Unexpected memory range for GPU %d", i)
	}
	t.Log("GPU devices initialized correctly")

	// group into a communicator ig
	rank2Address := make(map[uint32]string)
	for i := uint32(0); i < uint32(numGpus); i++ {
		addr := "127.0.0.1:500" + fmt.Sprint(i+1)
		rank2Address[i] = addr
	}

	for i, gpu := range gpus {
		_, err := gpu.SetupCommunication(
			context.Background(),
			&pb.SetupCommunicationRequest{
				RankToAddress: rank2Address,
				Rank:          &pb.Rank{Value: uint32(i)},
			})
		if err != nil {
			t.Logf("Error in SetupCommunication: %v", err)
		}
		assert.Equal(t, nil, err)
	}

	// perform send and receive
	sendResp, err := gpus[1].BeginSend(context.Background(),
		&pb.BeginSendRequest{
			SendBuffAddr: &pb.MemAddr{Value: gpuMetadatas[1].MinMemAddr},
			NumBytes:     0x100,
			DstRank:      &pb.Rank{Value: 0},
		})
	assert.Equal(t, nil, err)
	t.Logf("streamid: %d", sendResp.StreamId.Value)

	// why are these fields optional, is there a way to make them not?
	_, err = gpus[0].BeginReceive(context.Background(),
		&pb.BeginReceiveRequest{
			RecvBuffAddr: &pb.MemAddr{Value: gpuMetadatas[0].MinMemAddr},
			NumBytes:     0x100,
			StreamId:     sendResp.StreamId,
			SrcRank:      &pb.Rank{Value: 1},
			Op:           pb.ReduceOp_SUM,
		})
	assert.Equal(t, nil, err)

	status_response, err := gpus[0].GetStreamStatus(context.Background(),
		&pb.GetStreamStatusRequest{
			StreamId: sendResp.StreamId,
			SrcRank:  &pb.Rank{Value: 1},
		})
	// assert.Equal(t, pb.Status_IN_PROGRESS, status_response.Status)
	assert.Equal(t, nil, err)
	t.Logf("status resp: %v", status_response)
}
