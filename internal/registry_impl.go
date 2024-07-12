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
func (r *OrionRegistryImplementations) SubscribeToStream(a *proto.InitializeRequest, z proto.Registry_SubscribeToStreamServer) error {
	// check if the client is authorized
	// TODO: We should verify that the data is signed using the certificate
	// TODO: and that the certificate is valid under the CA and allowed to
	// TODO: use this user id.
	data := fmt.Sprintf("%d:%s:%d", a.MemberId, a.FriendlyName, a.TimestampSigned)
	block, _ := pem.Decode(a.Certificate)

	certs, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		log.Debug().Err(err).Msg("failed to parse certificate")
		return err
	}
	cert := certs[0]
	hash := sha512.New().Sum([]byte(data))
	successful := ecdsa.VerifyASN1(cert.PublicKey.(*ecdsa.PublicKey), hash, a.Signed)
	if successful {
		log.Debug().Err(err).Msg("signature is not valid")
		return fmt.Errorf("signature does not seem to be valid")
	}

	root := x509.NewCertPool()
	ca, err := os.ReadFile("ca/ca.crt")
	if err != nil {
		log.Debug().Err(err).Msg("failed to load CA certificate")
		return fmt.Errorf("signature does not seem to be valid")
	}
	ok := root.AppendCertsFromPEM(ca)
	if !ok {
		panic("failed to add root cert")
	}

	intermediates := x509.NewCertPool()
	ok = intermediates.AppendCertsFromPEM(a.Certificate)
	if !ok {
		panic("failed to add root cert")
	}
	opts := x509.VerifyOptions{
		Roots:         root,
		Intermediates: intermediates,
		DNSName:       fmt.Sprintf("%d.mem.orionet.re", a.MemberId),
	}

	if _, err := cert.Verify(opts); err != nil {
		log.Debug().Err(err).Msg("certificate is not valid for wanted domains")
		return fmt.Errorf(err.Error())
	}

	log.Info().Msgf("User %s auth with %s", opts.DNSName, cert.SerialNumber)
	// all good
	return nil
}

// When an existing client wants to initiate a connection to a new or existing peer.
func (r *OrionHolePunchingImplementation) InitializeConnectionToPeer(context.Context, *proto.InitializeConnectionToPeerRequest) (*proto.InitializeConnectionToPeerResponse, error) {
	return nil, nil
}
