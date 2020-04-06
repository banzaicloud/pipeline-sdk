package http

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"emperror.dev/errors"

	"github.com/banzaicloud/pipeline-sdk/.gen/pipeline/pipeline"
	"github.com/banzaicloud/pipeline-sdk/transport"
)

type Transport struct {
	configuration *pipeline.Configuration
}

func (t *Transport) Configuration() *pipeline.Configuration {
	return t.configuration
}

func NewTransport(config transport.Config) (*Transport, error) {
	rootCAs := x509.NewCertPool()

	caCertPEM, err := ioutil.ReadFile(config.CACertFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read HTTP CA certificate")
	}

	if !rootCAs.AppendCertsFromPEM(caCertPEM) {
		return nil, errors.New("failed to append HTTP CA certificate")
	}

	clientCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load HTTP client certificate")
	}

	tlsConfig := tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{clientCert},
	}

	transport := http.Transport{TLSClientConfig: &tlsConfig}
	client := http.Client{Transport: &transport}

	configuration := pipeline.NewConfiguration()

	configuration.HTTPClient = &client
	configuration.Scheme = "https"
	configuration.Host = config.Address
	configuration.BasePath = "pipeline"

	return &Transport{configuration: configuration}, nil
}
