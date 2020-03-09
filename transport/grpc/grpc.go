package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"emperror.dev/errors"
	"github.com/banzaicloud/pipeline-internal-sdk/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Transport struct {
	clientConn *grpc.ClientConn
}

func (t *Transport) ClientConn() *grpc.ClientConn {
	return t.clientConn
}

func NewTransport(config transport.Config) (*Transport, error) {
	rootCAs := x509.NewCertPool()

	caCertPEM, err := ioutil.ReadFile(config.CACertFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read GRPC CA certificate")
	}

	if !rootCAs.AppendCertsFromPEM(caCertPEM) {
		return nil, errors.New("failed to append GRPC CA certificate")
	}

	clientCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load GRPC client certificate")
	}

	tlsConfig := tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{clientCert},
	}

	creds := credentials.NewTLS(&tlsConfig)

	conn, err := grpc.Dial(config.Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	return &Transport{clientConn: conn}, nil
}
