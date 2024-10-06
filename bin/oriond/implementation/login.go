package implementation

import (
	"context"
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
	log.Info().Msg("loading the certificate")

	// Load the certificates
	certificateFile, err := os.ReadFile(*internal.CertificatePath)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", *internal.CertificatePath).
			Msg("cannot open the certificate path")
		return err
	}

	keyFile, err := os.ReadFile(*internal.KeyPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", *internal.KeyPath).
			Msg("cannot open the key path")
		return err
	}

	// Parsing the PEM file as a PKCS8 private key
	rawCertificate, _ := pem.Decode(keyFile)
	pk, err := x509.ParsePKCS8PrivateKey(rawCertificate.Bytes)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", *internal.KeyPath).
			Msg("failed to parse the pkcs8-formated private key")
		return err
	}

	// We check the key type to match ecdsa
	ecdsaKey, isEcdsaKey := pk.(*ecdsa.PrivateKey)
	if !isEcdsaKey {
		err = fmt.Errorf("this private key is not a ECDSA private key, oriond only works with ECDSA private keys")
		log.Error().
			Err(err).
			Str("file", *internal.KeyPath).
			Msg(err.Error())
		return err
	}

	identity := proto.Router{}

	// We compute the nonce with all the data
	nonce, err := internal.CalculateNonce(&identity, certificateFile, ecdsaKey)
	if err != nil {
		err = fmt.Errorf("couldn't compute the nonce given from the given information about this node")
		log.Error().
			Err(err).
			Msg(err.Error())
		return err
	}

	if c.sID != "" {
		nonce.SessionId = c.sID
		nonce.Reconnect = true
	}

	// We send the login-initialize information to the server
	if err = c.registryStream.Send(&proto.PeersToServer{
		Event: &proto.PeersToServer_Login{
			// TODO: More data for login
			Login: &proto.RouterLogin{},
		},
	}); err != nil {
		log.Error().
			Err(err).
			Msg("couldn't send the login signature to the server")
	}

	return nil
}

func (c *OrionClientDaemon) Start() error {
	if c.ctxCancel != nil && c.Context != nil {
		c.ctxCancel()
	}

	ctx, cancel := context.WithCancel(c.ParentCtx)
	c.ctxCancel = cancel
	c.Context = ctx

	// Intialize the streams to the signaling server
	if err := c.initializeStream(); err != nil {
		return err
	}
	if err := c.login(); err != nil {
		return err
	}

	// Start the listener as a background task
	go c.listener()
	return nil
}
