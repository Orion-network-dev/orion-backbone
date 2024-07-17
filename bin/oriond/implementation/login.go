package implementation

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (c *OrionClientDaemon) login() error {
	log.Info().Msg("Loading the certificates")

	// Load the certificates
	certificateFile, err := os.ReadFile(*internal.CertificatePath)
	if err != nil {
		log.Error().Err(err).Str("file", *internal.CertificatePath).Msgf("Cannot open the certificate path")
		return err
	}

	// Loading the private key data
	privateKey, err := os.ReadFile(*internal.KeyPath)
	if err != nil {
		log.Error().Err(err).Msg("this server's private key couldn't be read")
		return err
	}

	// Parsing the PEM file as a PKCS8 private key
	rawCertificate, _ := pem.Decode(privateKey)
	pk, err := x509.ParsePKCS8PrivateKey(rawCertificate.Bytes)
	if err != nil {
		log.Error().Err(err).Msgf("coundn't read the certificate key file")
		return err
	}

	// Once the private key is loaded, we check if it's a EC type key.
	ecdsaKey, isEcdsaKey := pk.(*ecdsa.PrivateKey)
	if !isEcdsaKey {
		err = fmt.Errorf("this private key is not an ECdsa-based key and cannot be used for Orion currently")
		log.Error().Err(err).Msg(err.Error())
		return err
	}
	nonce, err := internal.CalculateNonce(c.memberId, *friendlyName, certificateFile, ecdsaKey)
	if err != nil {
		log.Error().Err(err).Msgf("coundn't compute the nonce")
		return err
	}
	err = c.registryStream.Send(&proto.RPCClientEvent{
		Event: &proto.RPCClientEvent_Initialize{
			Initialize: nonce,
		},
	})
	if err != nil {
		log.Error().Err(err).Msgf("couldn't initialize the login to the grpc connection")
	}

	return nil
}
