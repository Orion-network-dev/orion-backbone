package implementation

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"software.sslmate.com/src/go-pkcs12"
)

func (c *OrionClientDaemon) login() error {
	log.Info().Msg("loading the certificate")
	p12, err := os.ReadFile(*internal.AuthorityPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to load the credentials file")
		return err
	}
	password, err := os.ReadFile(*internal.PasswordFile)
	if err != nil {
		log.Error().Err(err).Msg("failed to load the password file")
		return err
	}

	pk, certificate, _, err := pkcs12.DecodeChain(p12, string(password))
	if err != nil {
		log.Error().Err(err).Msg("failed to use the p12 file")
		return err
	}

	// We check the key type to match ecdsa
	ecdsaKey, isEcdsaKey := pk.(*ecdsa.PrivateKey)
	if !isEcdsaKey {
		err = fmt.Errorf("this private key is not a ECDSA private key, oriond only works with ECDSA private keys")
		log.Error().
			Err(err).
			Msg(err.Error())
		return err
	}

	publicKeyDer, _ := x509.MarshalPKIXPublicKey(certificate)

	// We compute the nonce with all the data
	nonce, err := internal.CalculateNonce(c.memberId, *friendlyName, publicKeyDer, ecdsaKey)
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
	if err = c.registryStream.Send(&proto.RPCClientEvent{
		Event: &proto.RPCClientEvent_Initialize{
			Initialize: nonce,
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
