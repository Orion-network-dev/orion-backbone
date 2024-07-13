package internal

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

var (
	AuthorityPath   = flag.String("tls-authority-path", "ca/ca.crt", "Path to the certificate authority file")
	CertificatePath = flag.String("tls-certificate-path", "", "Path to the certificate authority file")
	KeyPath         = flag.String("tls-key-path", "", "Path to the certificate authority file")
)

func loadAuthorityPool() (*x509.CertPool, error) {
	// Load the CA certificate
	trustedCert, err := os.ReadFile(*AuthorityPath)
	if err != nil {
		return nil, err
	}

	// Put the CA certificate into the certificate pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCert) {
		return nil, fmt.Errorf("the CA certificate couldn't be constructed")
	}

	return certPool, nil
}

func LoadTLS(clientCerts bool) (credentials.TransportCredentials, error) {
	// Load the client certificate and its key
	clientCert, err := tls.LoadX509KeyPair(*CertificatePath, *KeyPath)
	if err != nil {
		return nil, err
	}
	authorityPool, err := loadAuthorityPool()
	if err != nil {
		return nil, err
	}

	// Create the TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      authorityPool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	if clientCerts {
		tlsConfig.ClientCAs = authorityPool
	}

	// Return new TLS credentials based on the TLS configuration
	return credentials.NewTLS(tlsConfig), nil
}
