package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sl "github.com/JoshuaMBa/dsml/gpu_sim/gpu_device/server_lib"
	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
)

func TestServerBasic(t *testing.T) {

	gpu, _ := sl.MakeGPUDeviceServer(*sl.DefaultGPUDeviceOptions())
	go gpu.StreamSendDispatch()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := gpu.GetDeviceMetadata(ctx, &pb.GetDeviceMetadataRequest{})
	assert.True(t, err == nil)
	assert.True(t, out.Metadata.DeviceId.Value == 11)

}
