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
func calculateNonceBytes(MemberId int64, FriendlyName string, time int64) []byte {
	return sha512.New().Sum([]byte(fmt.Sprintf("%d:%s:%d", MemberId, FriendlyName, time)))
}

// Used to calculate the nonce and sign it
func CalculateNonce(
	MemberId int64,
	FriendlyName string,
	Certificate []byte,
	PrivateKey *ecdsa.PrivateKey,
) *proto.InitializeRequest {
	time := time.Now().Unix()
	authHash := calculateNonceBytes(MemberId, FriendlyName, time)

	signed, err := ecdsa.SignASN1(rand.Reader, PrivateKey, authHash)

	if err != nil {
		log.Fatal().Err(err).Msgf("couldn't sign the nonce data")
	}

	return &proto.InitializeRequest{
		FriendlyName:    FriendlyName,
		TimestampSigned: time,
		MemberId:        MemberId,
		Certificate:     Certificate,
		Signed:          signed,
	}
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

func Authenticate(
	TimeStamp int64,
	Certificate []byte,
	SignedNonce []byte,
	MemberId int64,
	FriendlyName string,
	RootCertPool *x509.CertPool,
) error {
	// Verify that the date only has a variation inferior to 2s
	time := time.Now().Unix()
	if time-TimeStamp > 2 {
		err := fmt.Errorf("the verification timestamp is too far from the current time")
		log.Debug().Err(err).Msg("user supplied an invalid date/time")
		return err
	}

	// Parse the user-given certificate
	cert, err := ParsePEMCertificate(Certificate)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse the user leaf certificate")
		return err
	}

	// Create a new pool from the user-given PEM trust chain
	intermediates, err := CreateCertPoolFromPEM(Certificate)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse the intermediary certificates")
		return err
	}

	identifier := fmt.Sprintf("%d.member.orionet.re", MemberId)

	// Verifying the certificate validity using the root certificate and user-provided
	// intermediary certificates. This checks that the certificate is signed and allowed to use
	// the name `{member_id}.mem.orionet.re` which specifies a member member for the member_id {member_id}
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:         RootCertPool,
		Intermediates: intermediates,
		DNSName:       identifier,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}); err != nil {
		log.Debug().Err(err).Msg("certificate is not valid for wanted domains")
		return fmt.Errorf(err.Error())
	}

	// Calculate the hash given in order to check the client signature
	nonce := calculateNonceBytes(MemberId, FriendlyName, TimeStamp)

	// Verify that the user-provided data matches the signature created using the client root key
	successful := ecdsa.VerifyASN1(cert.PublicKey.(*ecdsa.PublicKey), nonce, SignedNonce)
	if !successful {
		err := fmt.Errorf("this signature does not seem to be a valid ECDSA signature")
		log.Debug().Err(err).Msg("the user authentication failed, invalid signature")
		return err
	}

	log.Info().Msgf("User %s auth with certificate with serial: %s", identifier, cert.SerialNumber)

	return nil
}
