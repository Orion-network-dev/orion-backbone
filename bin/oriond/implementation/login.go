package implementation

import (
	"bytes"
	"context"
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
	certificateFile, err := os.ReadFile(*internal.TLSPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", *internal.TLSPath).
			Msg("cannot open the certificate path")
		return err
	}

	var buffer bytes.Buffer

	for block, rest := pem.Decode(certificateFile); block != nil; block, rest = pem.Decode(rest) {
		if block.Type == "CERTIFICATE" {
			err := pem.Encode(&buffer, block)
			if err != nil {
				log.Error().
					Err(err).
					Msg("cannot encode to pem")
				return err
			}
		}
	}

	// We compute the nonce with all the data
	nonce, err := internal.CalculateNonce(c.memberId, *friendlyName, buffer.Bytes(), c.privateKey)
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
