package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	pb "github.com/JoshuaMBa/dsml/gpu_sim/proto"
	sl "github.com/JoshuaMBa/dsml/gpu_sim/gpu_device/server_lib"
)

func TestServerBasic(t *testing.T) {

	gpuDeviceService, _ := sl.MakeGPUDeviceServer(*sl.DefaultGPUDeviceOptions())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := gpuDeviceService.GetDeviceMetadata(ctx, &pb.GetDeviceMetadataRequest{})
	assert.True(t, err == nil)
	assert.True(t, out.Metadata.DeviceId.Value == 11)
}
