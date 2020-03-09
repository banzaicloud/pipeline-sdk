package main

import (
	"context"
	"time"

	"github.com/banzaicloud/pipeline-internal-sdk/.gen/pipeline/pipeline"
	"github.com/banzaicloud/pipeline-internal-sdk/process"
	"github.com/banzaicloud/pipeline-internal-sdk/transport"
	"github.com/banzaicloud/pipeline-internal-sdk/transport/grpc"
	"github.com/banzaicloud/pipeline-internal-sdk/transport/http"
)

const processID = "1234-5678-9111-1234"

func main() {
	config := transport.Config{
		Address:    "127.0.0.1:9092",
		CACertFile: "config/certs/ca.pem",
		CertFile:   "config/certs/client.pem",
		KeyFile:    "config/certs/client-key.pem",
	}

	grpcTransport, err := grpc.NewTransport(config)
	if err != nil {
		panic(err)
	}

	c, err := process.NewClient(grpcTransport)
	if err != nil {
		panic(err)
	}

	p := process.ProcessEntry{
		ID:           processID,
		OrgID:        1,
		Name:         "cluster-create",
		Status:       process.RunningStatus,
		ResourceType: process.ClusterResourceType,
		ResourceID:   "13",
		StartedAt:    time.Now(),
	}

	err = c.LogProcess(context.Background(), p)
	if err != nil {
		panic(err)
	}

	{
		err = c.LogEvent(context.Background(), process.ProcessEvent{
			ProcessID: processID,
			Name:      "create-vpc",
			Log:       "Creating VPC in AWS",
			Timestamp: time.Now(),
		})
		if err != nil {
			panic(err)
		}

		err = c.LogEvent(context.Background(), process.ProcessEvent{
			ProcessID: processID,
			Name:      "create-vpc",
			Log:       "Creating VPC in AWS finished",
			Timestamp: time.Now(),
		})
		if err != nil {
			panic(err)
		}
	}

	finishedAt := time.Now()
	p.Status = process.FinishedStatus
	p.FinishedAt = &finishedAt

	err = c.LogProcess(context.Background(), p)
	if err != nil {
		panic(err)
	}

	////
	config.Address = "127.0.0.1:9090"

	httpTransport, err := http.NewTransport(config)
	if err != nil {
		panic(err)
	}

	apiClient := pipeline.NewAPIClient(httpTransport.Configuration())

	secrets, _, err := apiClient.SecretsApi.GetSecrets(context.Background(), 1, nil)
	if err != nil {
		panic(err)
	}

	for _, secret := range secrets {
		println(secret.Name)
	}
}
