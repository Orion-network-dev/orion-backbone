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
		case *proto.RPCServerEvent_NewMember:
			c.handleNewClient(subCtx, event.Event.(*proto.RPCServerEvent_NewMember).NewMember)
		case *proto.RPCServerEvent_MemberConnect:
			c.handleWantsToConnect(subCtx, event.Event.(*proto.RPCServerEvent_MemberConnect).MemberConnect)
		case *proto.RPCServerEvent_DisconnectedMember:
			c.handleRemovedClient(event.Event.(*proto.RPCServerEvent_DisconnectedMember).DisconnectedMember)
		case *proto.RPCServerEvent_MemberConnectResponse:
			c.handleWantsToConnectResponse(event.Event.(*proto.RPCServerEvent_MemberConnectResponse).MemberConnectResponse)
		case *proto.RPCServerEvent_SessionId:
			c.sID = event.Event.(*proto.RPCServerEvent_SessionId).SessionId
		}
		cancel()
	}
}
