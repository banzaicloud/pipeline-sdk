package main

import (
	"context"

	"github.com/banzaicloud/pipeline-sdk/.gen/pipeline/pipeline"
	"github.com/banzaicloud/pipeline-sdk/transport"
	"github.com/banzaicloud/pipeline-sdk/transport/http"
)

const processID = "1234-5678-9111-1234"

func main() {
	config := transport.Config{
		Address:    "127.0.0.1:9090",
		CACertFile: "config/certs/ca.pem",
		CertFile:   "config/certs/client.pem",
		KeyFile:    "config/certs/client-key.pem",
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
