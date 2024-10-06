package implementation

import (
	"context"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/link"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (c *OrionClientDaemon) handleNewRouter(
	ctx context.Context,
	event *proto.RouterConnectedEvent,
) {
	c.tunnelsLock.Lock()
	defer c.tunnelsLock.Unlock()
	peerID := internal.IdentityFromRouter(event.Router)
	log.Info().
		Uint64("peer-id", peerID).
		Msg("Initiating a new link")

	if c.tunnels[peerID] != nil {
		log.Error().
			Uint64("peer-id", peerID).
			Msg("received a new_client event for a already-initialized event")
		return
	}

	peer, err := link.NewPeerLink(
		c.Context,
		c.identity,
		event.Router,
		c.wgClient,
		c.frrManager,
	)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("peer-id", peerID).
			Msgf("failed to initialize the peer object")
		return
	}
	// We save the peer in the tunnels map
	c.tunnels[peerID] = peer
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
			Uint64("peer-id", peerID).
			Msg("failed to hole-punch the interface")
		peer.Dispose()
		return
	}

	// Inform the peer that we re ready for connection
	err = c.registryStream.Send(&proto.PeersToServer{
		Event: &proto.PeersToServer_Initiate{
			Initiate: &proto.RouterPeerToPeerInitiate{
				EndpointAddr: holePunching.ClientEndpointAddr,
				EndpointPort: holePunching.ClientEndpointPort,
				PublicKey:    publickey[:],
				FriendlyName: *friendlyName,
				Routing: &proto.RoutingInformation{
					Destination: event.Router,
					Source:      c.identity,
				},
			},
		},
	})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("peer-id", peerID).
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
				if establishedStreamID == peerID {
					return
				}
			case <-waitingForResponse.Done():
				// timeout reached while establishing the peer connection
				log.Error().
					Err(waitingForResponse.Err()).
					Uint64("peer-id", peerID).
					Msgf("timeout exceeded while waiting for a response from the peer")
				peer.Dispose()
				return
			}
		}
	}()
}
