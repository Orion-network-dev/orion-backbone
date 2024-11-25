package internal

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/pkcs12"
	"google.golang.org/grpc/credentials"
)

var (
	AuthorityPath = flag.String("tls-authority-path", "/etc/oriond/ca.crt", "Path to the certificate authority file")
	P12Path       = flag.String("tls-authentication", "/etc/oriond/identity.p12", "Path to the p12 file")
	PasswordFile  = flag.String("tls-authentication-pass-file", "", "Password for the identity")
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
	log.Debug().Str("authority-path", *AuthorityPath).Str("certificate-path", *P12Path).Str("pass-path", *PasswordFile).Msg("loading the certificates for login")

	p12, err := os.ReadFile(*AuthorityPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to load the credentials file")
		return nil, err
	}
	password, err := os.ReadFile(*PasswordFile)
	if err != nil {
		log.Error().Err(err).Msg("failed to load the password file")
		return nil, err
	}

	key, cert, err := pkcs12.Decode(p12, string(password))
	if err != nil {
		log.Error().Err(err).Msg("failed to use the p12 file")
		return nil, err
	}

	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key.(crypto.PrivateKey),
		Leaf:        cert,
	}

	authorityPool, err := loadAuthorityPool()
	if err != nil {
		log.Error().Err(err).Msg("failed to load the authority files")
		return nil, err
	}

	// Create the TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
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
