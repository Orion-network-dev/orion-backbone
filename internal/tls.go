package internal

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/credentials"
)

var (
	AuthorityPath   = flag.String("tls-authority-path", "/etc/oriond/ca.crt", "Path to the certificate authority file")
	CertificatePath = flag.String("tls-certificate-path", "/etc/oriond/identity.crt", "Path to the certificate authority file")
	KeyPath         = flag.String("tls-key-path", "/etc/oriond/identity.key", "Path to the certificate authority file")
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
	log.Debug().Str("authority-path", *AuthorityPath).Str("certificate-path", *CertificatePath).Str("key-path", *KeyPath).Msg("loading the certificates for login")

	// Load the client certificate and its key
	clientCert, err := tls.LoadX509KeyPair(*CertificatePath, *KeyPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to load the certificate")
		return nil, err
	}
	authorityPool, err := loadAuthorityPool()
	if err != nil {
		log.Error().Err(err).Msg("failed to load the authority files")
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
