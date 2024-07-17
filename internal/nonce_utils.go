package internal

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

// Calculate the nonce bytes
func CalculateNonceBytes(MemberId uint32, FriendlyName string, time int64) []byte {
	return sha512.New().Sum([]byte(fmt.Sprintf("%d:%s:%d", MemberId, FriendlyName, time)))
}

// Used to calculate the nonce and sign it
func CalculateNonce(
	MemberId uint32,
	FriendlyName string,
	Certificate []byte,
	PrivateKey *ecdsa.PrivateKey,
) (*proto.InitializeRequest, error) {
	time := time.Now().Unix()
	authHash := CalculateNonceBytes(MemberId, FriendlyName, time)

	signed, err := ecdsa.SignASN1(rand.Reader, PrivateKey, authHash)

	if err != nil {
		log.Error().Err(err).Msgf("couldn't sign the nonce data")
		return nil, err
	}

	return &proto.InitializeRequest{
		FriendlyName:    FriendlyName,
		TimestampSigned: time,
		MemberId:        MemberId,
		Certificate:     Certificate,
		Signed:          signed,
	}, nil
}

// Parse a pem file and adds all the pem-encoded certificates to a cert-pool
func CreateCertPoolFromPEM(PEMData []byte) (*x509.CertPool, error) {
	// Create a cert-pool containing the user-provided intermediary certificates
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(PEMData)
	if !ok {
		user_err := fmt.Errorf("failed to import the User-given intermediate CAs")
		log.Debug().Err(user_err).Msg("failed to import the root ca certificate")
		return nil, user_err
	}
	return pool, nil
}

func ParsePEMCertificate(Certificate []byte) (*x509.Certificate, error) {
	// Parsing the pem-encoded used-given certificate in order to parse the certificate
	block, _ := pem.Decode(Certificate)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse certificate")
		return nil, err
	}
	return cert, nil
}
