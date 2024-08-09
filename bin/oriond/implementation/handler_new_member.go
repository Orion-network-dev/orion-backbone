package implementation

import (
	"context"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/link"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (c *OrionClientDaemon) handleNewClient(
	ctx context.Context,
	event *proto.NewMemberEvent,
) {
	c.tunnelsLock.Lock()
	defer c.tunnelsLock.Unlock()

	log.Info().
		Uint32("peer-id", event.PeerId).
		Msg("Initiating a new link")

	if c.tunnels[event.PeerId] != nil {
		log.Error().
			Uint32("peer-id", event.PeerId).
			Msg("received a new_client event for a already-initialized event")
		return
	}

	peer, err := link.NewPeerLink(
		c.Context,
		c.memberId,
		event.PeerId,
		c.wgClient,
		c.frrManager,
	)
	if err != nil {
		log.Error().
			Err(err).
			Uint32("peer-id", event.PeerId).
			Msgf("failed to initialize the peer object")
		return
	}
	// We save the peer in the tunnels map
	c.tunnels[event.PeerId] = peer
	publickey := peer.PublicKey()

	// We initialize a one minte context for getting the hole-punching details
	holePunchingContext, cancel := context.WithTimeout(ctx, time.Minute)
	holePunching, err := peer.HolePunchTunnel(
		holePunchingContext,
		c.holePunchingClient,
	)
	cancel()
	if err != nil {
		log.Error().
			Err(err).
			Uint32("peer-id", event.PeerId).
			Msg("failed to hole-punch the interface")
		peer.Dispose()
		return
	}

	// Inform the peer that we re ready for connection
	err = c.registryStream.Send(&proto.RPCClientEvent{
		Event: &proto.RPCClientEvent_Connect{
			Connect: &proto.MemberConnectEvent{
				EndpointAddr:      holePunching.ClientEndpointAddr,
				EndpointPort:      holePunching.ClientEndpointPort,
				PublicKey:         publickey[:],
				FriendlyName:      *friendlyName,
				DestinationPeerId: event.PeerId,
				SourcePeerId:      c.memberId,
			},
		},
	})
	if err != nil {
		log.Error().
			Err(err).
			Uint32("peer-id", event.PeerId).
			Msgf("couldn't write the initialization message to the gRPC connection")
		peer.Dispose()
		return
	}
	go func() {
		waitingForResponse, cancel := context.WithTimeout(c.Context, time.Minute)
		defer cancel()
		establshed_stream := c.establishedStream.Listener(10)
		defer establshed_stream.Close()
		for {
			select {
			case establishedStreamID := <-establshed_stream.Ch():
				// Our connection got succesfully established
				if establishedStreamID == event.PeerId {
					return
				}
			case <-waitingForResponse.Done():
				// timeout reached while establishing the peer connection
				log.Error().
					Err(waitingForResponse.Err()).
					Uint32("peer-id", event.PeerId).
					Msgf("timeout exceeded while waiting for a response from the peer")
				peer.Dispose()
				return
			}
		}
	}()
}
