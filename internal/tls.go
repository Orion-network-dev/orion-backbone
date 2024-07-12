package internal

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"google.golang.org/grpc/credentials"
)

func NewClientTLS() credentials.TransportCredentials {
	// Load the client certificate and its key
	clientCert, err := tls.LoadX509KeyPair("ca/client.crt", "ca/client.key")

	if err != nil {
		log.Fatalf("Failed to load client certificate and key %v", err)
	}

	// Load the CA certificate
	trustedCert, err := os.ReadFile("ca/ca.crt")
	if err != nil {
		log.Fatalf("Failed to load trusted certificate %v", err)
	}

	// Put the CA certificate into the certificate pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCert) {
		log.Fatalf("Failed to append trusted certificate to certificate pool %v", err)
	}

	// Create the TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	// Return new TLS credentials based on the TLS configuration
	return credentials.NewTLS(tlsConfig)
}

func NewServerTLS() credentials.TransportCredentials {
	// Load the server certificate and its key
	serverCert, err := tls.LoadX509KeyPair("ca/server.crt", "ca/server.key")

	if err != nil {
		log.Fatalf("Failed to load server certificate and key %v", err)
	}

	// Load the CA certificate
	trustedCert, err := os.ReadFile("ca/ca.crt")
	if err != nil {
		log.Fatalf("Failed to load trusted certificate %v", err)
	}

	// Put the CA certificate into the certificate pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCert) {
		log.Fatalf("Failed to append trusted certificate to certificate pool %v", err)
	}

	// Create the TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		RootCAs:      certPool,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	// Return new TLS credentials based on the TLS configuration
	return credentials.NewTLS(tlsConfig)
}
