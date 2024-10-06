package implementation

import (
	"context"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (c *OrionClientDaemon) initializeStream() error {
	ctx := c.Context

	stream, err := c.registryClient.SubscribeToStream(ctx)
	if err != nil {
		return err
	}
	c.registryStream = stream
	return nil
}

func (c *OrionClientDaemon) listener() error {
	defer func() {
		log.Info().Msg("listener is finished")
		c.ctxCancel()
	}()

	for {
		event, err := c.registryStream.Recv()
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to read the stream from the registry")
			return err
		}

		subCtx, cancel := context.WithTimeout(c.Context, time.Second*10)
		switch event.Event.(type) {
		case *proto.ServerToPeers_Connected:
			c.handleNewRouter(subCtx, event.Event.(*proto.ServerToPeers_Connected).Connected)
		case *proto.ServerToPeers_Disconnected:
			c.handleDisconnectedRouter(subCtx, event.Event.(*proto.ServerToPeers_Disconnected).Disconnected)
		case *proto.ServerToPeers_Initiate:
			c.handleInitiate(subCtx, event.Event.(*proto.ServerToPeers_Initiate).Initiate)
		case *proto.ServerToPeers_InitiateAck:
			c.handleInitiateAck(subCtx, event.Event.(*proto.ServerToPeers_InitiateAck).InitiateAck)
		case *proto.ServerToPeers_Terminate:
			c.handleTerminate(subCtx, event.Event.(*proto.ServerToPeers_Terminate).Terminate)
		case *proto.ServerToPeers_TerminateAck:
			c.handleTerminateAck(subCtx, event.Event.(*proto.ServerToPeers_TerminateAck).TerminateAck)
		case *proto.ServerToPeers_SessionId:
			log.Info().
				Msg("got a sessionId from the registry server")
			c.sID = event.Event.(*proto.ServerToPeers_SessionId).SessionId.SessionId
		}
		cancel()
	}
}
