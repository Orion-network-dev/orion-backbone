package implementation

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
)

func (c *OrionClientDaemon) initializeStream() error {
	ctx := c.ctx

	stream, err := c.registryClient.SubscribeToStream(ctx)
	if err != nil {
		return err
	}
	c.registryStream = stream
	return nil
}

func (c *OrionClientDaemon) listener() error {
	for {
		event, err := c.registryStream.Recv()
		if err != nil {
			return err
		}

		switch event.Event.(type) {
		case *proto.RPCServerEvent_NewClient:
			c.handleNewClient(event.Event.(*proto.RPCServerEvent_NewClient).NewClient)
		case *proto.RPCServerEvent_RemovedClient:
			c.handleRemovedClient(event.Event.(*proto.RPCServerEvent_RemovedClient).RemovedClient)
		case *proto.RPCServerEvent_WantsToConnect:
			c.handleWantsToConnect(event.Event.(*proto.RPCServerEvent_WantsToConnect).WantsToConnect)
		case *proto.RPCServerEvent_WantsToConnectResponse:
			c.handleWantsToConnectResponse(event.Event.(*proto.RPCServerEvent_WantsToConnectResponse).WantsToConnectResponse)
		}
	}
}
