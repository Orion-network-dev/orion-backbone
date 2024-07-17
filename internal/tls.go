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
	// AuthorityPath is the path to the certificate authority file.
	AuthorityPath = flag.String("tls-authority-path", "/etc/oriond/ca.crt", "Path to the certificate authority file")
	// CertificatePath is the path to the client certificate file.
	CertificatePath = flag.String("tls-certificate-path", "/etc/oriond/identity.crt", "Path to the client certificate file")
	// KeyPath is the path to the client key file.
	KeyPath = flag.String("tls-key-path", "/etc/oriond/identity.key", "Path to the client key file")
)

// LoadTLS loads the TLS configuration using the client certificate, key, and CA certificate.
// If clientCerts is true, it also sets up client authentication.
func LoadTLS(clientCerts bool) (credentials.TransportCredentials, error) {
	log.Debug().
		Str("authority-path", *AuthorityPath).
		Str("certificate-path", *CertificatePath).
		Str("key-path", *KeyPath).
		Msg("loading the certificates for login")

	clientCert, err := tls.LoadX509KeyPair(*CertificatePath, *KeyPath)
	if err != nil {
		return nil, HandleError(err, "failed to load client certificate")
	}

	authorityPool, err := loadAuthorityPool()
	if err != nil {
		return nil, HandleError(err, "failed to load the authority pool")
	}

	tlsConfig := createTLSConfig(clientCert, authorityPool, clientCerts)

	return credentials.NewTLS(tlsConfig), nil
}

// createTLSConfig creates a TLS configuration using the provided client certificate and CA pool.
// If clientCerts is true, it also sets up client authentication with the CA pool.
func createTLSConfig(clientCert tls.Certificate, authorityPool *x509.CertPool, clientCerts bool) *tls.Config {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      authorityPool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	if clientCerts {
		tlsConfig.ClientCAs = authorityPool
	}

	return tlsConfig
}

// loadAuthorityPool loads the CA certificate into a new x509.CertPool.
func loadAuthorityPool() (*x509.CertPool, error) {
	trustedCert, err := os.ReadFile(*AuthorityPath)
	if err != nil {
		return nil, HandleError(err, "failed to read CA certificate")
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCert) {
		return nil, HandleError(fmt.Errorf("the CA certificate couldn't be constructed"), "failed to append CA certificate to pool")
	}

	return certPool, nil
}
