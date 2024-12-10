package server_lib

import (
	"context"
	"log"
	"math/rand"

	"github.com/JoshuaMBa/dsml/failure_injection"
	fipb "github.com/JoshuaMBa/dsml/failure_injection/proto"
	gofakeit "github.com/brianvoe/gofakeit/v6"
)

// options for service init
type GPUDeviceOptions struct {
	// Random seed for generating database data
	Seed int64
	// failure injection config params
	SleepNs              int64
	FailureRate          int64
	ResponseOmissionRate int64
	// maximum size of batches accepted
	MaxBatchSize int
}

func DefaultUserServiceOptions() *GPUDeviceOptions {
	return &GPUDeviceOptions{
		Seed:                 42,
		SleepNs:              0,
		FailureRate:          0,
		ResponseOmissionRate: 0,
		MaxBatchSize:         50,
	}
}

type GPUDeviceServer struct {
	// options are read-only and intended to be immutable during the lifetime of
	// the service.  The failure injection config (in the FailureInjector) is
	// mutable (via SetInjectionConfigRequest-s) after the server starts.
	options GPUDeviceOptions
	fi      *failure_injection.FailureInjector
}

func MakeGPUDeviceServer(
	options GPUDeviceOptions,
) *GPUDeviceServer {
	rand.Seed(options.Seed)
	gofakeit.Seed(options.Seed)
	fi := failure_injection.MakeFailureInjector()
	fi.SetInjectionConfig(
		options.SleepNs,
		options.FailureRate,
		options.ResponseOmissionRate,
	)
	fiConfig := fi.GetInjectionConfig()
	log.Printf(
		"Starting UserService with failure injection config: [sleepNs: %d, failureRate: %d, responseOmissionRate: %d]",
		fiConfig.SleepNs,
		fiConfig.FailureRate,
		fiConfig.ResponseOmissionRate,
	)
	return &GPUDeviceServer{
		options: options,
		fi:      fi,
	}
}

func (db *GPUDeviceServer) SetInjectionConfig(
	ctx context.Context,
	req *fipb.SetInjectionConfigRequest,
) (*fipb.SetInjectionConfigResponse, error) {
	db.fi.SetInjectionConfigPb(req.Config)
	log.Printf(
		"UserService failure injection config set to: [sleepNs: %d, failureRate: %d, responseOmissionRate: %d]",
		req.Config.SleepNs,
		req.Config.FailureRate,
		req.Config.ResponseOmissionRate,
	)

	return &fipb.SetInjectionConfigResponse{}, nil
}
