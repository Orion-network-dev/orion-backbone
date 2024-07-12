package internal

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

type OrionRegistryImplementations struct {
	newPeers chan proto.InitializeRequest
	proto.UnimplementedRegistryServer
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// TODO: Implement p2p session creation through the registry.
func (r *OrionRegistryImplementations) SubscribeToStream(initializeRequest *proto.InitializeRequest, z proto.Registry_SubscribeToStreamServer) error {

	// Verifify the date
	time := time.Now().Unix()

	// Allow a time variation  of 2s
	if time-initializeRequest.TimestampSigned > 2 {
		err := fmt.Errorf("the verification timestamp is too far from the current time")
		log.Debug().Err(err).Msg("user supplied an invalid date/time")
		return err
	}

	// Calculate the hash using the user-given parameters.
	authString := fmt.Sprintf("%d:%s:%d", initializeRequest.MemberId, initializeRequest.FriendlyName, initializeRequest.TimestampSigned)
	authStringHash := sha512.New().Sum([]byte(authString))

	// Parse the user-given pem-certificate
	block, _ := pem.Decode(initializeRequest.Certificate)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse certificate")
		return err
	}

	// Verify the signature using ECDSA
	successful := ecdsa.VerifyASN1(cert.PublicKey.(*ecdsa.PublicKey), authStringHash, initializeRequest.Signed)
	if !successful {
		err := fmt.Errorf("this signature does not seem to be a valid ECDSA signature")
		log.Debug().Err(err).Msg("the user authentication failed, invalid signature")
		return err
	}

	// Reading the root certificate
	ca, err := os.ReadFile("ca/ca.crt")
	if err != nil {
		user_err := fmt.Errorf("failed to import the ROOT CA")
		log.Debug().Err(user_err).Err(err).Msg("failed to import the root ca certificate")
		return user_err
	}
	// Create a new certificate pool containing the root certificates
	root := x509.NewCertPool()
	// Append the root certificate to the pool
	ok := root.AppendCertsFromPEM(ca)
	if !ok {
		log.Debug().Msg("the root certificate couldn't be imported")
	}

	// Create a new certificate pool containing the user-provided intermadiary certificates
	intermediates := x509.NewCertPool()
	ok = intermediates.AppendCertsFromPEM(initializeRequest.Certificate)
	if !ok {
		user_err := fmt.Errorf("failed to import the User-given intermediate CAs")
		log.Debug().Err(user_err).Err(err).Msg("failed to import the root ca certificate")
		return user_err
	}

	// Handle the user certificate verification parameters
	opts := x509.VerifyOptions{
		Roots:         root,
		Intermediates: intermediates,
		DNSName:       fmt.Sprintf("%d.mem.orionet.re", initializeRequest.MemberId),
	}

	// Verify the certificate
	if _, err := cert.Verify(opts); err != nil {
		log.Debug().Err(err).Msg("certificate is not valid for wanted domains")
		return fmt.Errorf(err.Error())
	}

	// Our user is finally authenticated! Oorah!
	log.Info().Msgf("User %s auth with %s", opts.DNSName, cert.SerialNumber)

	return nil
}

// When an existing client wants to initiate a connection to a new or existing peer.
func (r *OrionHolePunchingImplementation) InitializeConnectionToPeer(context.Context, *proto.InitializeConnectionToPeerRequest) (*proto.InitializeConnectionToPeerResponse, error) {
	return nil, nil
}
