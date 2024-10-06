package implementation

import (
	"context"
	"net"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/link"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (c *OrionClientDaemon) handleInitiate(
	ctx context.Context,
	event *proto.RouterPeerToPeerInitiate,
) {

	source := internal.IdentityFromRouter(event.Routing.Source)
	destination := internal.IdentityFromRouter(event.Routing.Destination)
	self := internal.IdentityFromRouter(c.identity)

	// We ignore requests that are coming from ourself,
	if destination != self ||
		source == self {
		log.Error().
			Uint64("peer-id", source).
			Msg("received a message not destinated to this host")
		return
	}
	c.tunnelsLock.Lock()
	defer c.tunnelsLock.Unlock()
	if c.tunnels[source] != nil {
		log.Error().
			Uint64("peer-id", source).
			Msg("received a want to connect event for a already-initialized event")
		c.tunnels[source].Dispose()
	}

	// It's our job to generate the pre-shared key information
	// according the the protocol.
	preshared, err := wgtypes.GenerateKey()
	if err != nil {
		log.Error().
			Err(err).
			Uint64("peer-id", source).
			Msg("failed to send the message to the server")
		return
	}

	// We initialize a new peer link object to hold this new link context
	peer, err := link.NewPeerLink(
		c.Context,
		c.identity,
		event.Routing.Source,
		c.wgClient,
		c.frrManager,
	)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("peer-id", source).
			Msg("failed to initialize the peer link object")
		return
	}
	publicKey := peer.PublicKey()

	c.tunnels[source] = peer

	// We initialize a one minte context for getting the hole-punching details
	holePunchingContext, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	// We get the nat-ed tunnel details
	holePunching, err := peer.HolePunchTunnel(
		holePunchingContext,
		c.holePunchingClient,
	)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("peer-id", source).
			Msg("failed to hole-punch the interface")
		peer.Dispose()
		return
	}

	// We initialize the peer to configure the peer & bgp sessions
	err = peer.InitializePeerConnection(
		&net.UDPAddr{
			IP:   net.ParseIP(event.EndpointAddr),
			Port: int(event.EndpointPort),
		},
		wgtypes.Key(event.PublicKey),
		preshared,
	)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("peer-id", source).
			Msg("failed to initialize the connection to the peer")
		peer.Dispose()
		return
	}

	err = c.registryStream.Send(&proto.PeersToServer{
		Event: &proto.PeersToServer_InitiateAck{
			InitiateAck: &proto.RouterPeerToPeerInitiateACK{
				EndpointAddr: holePunching.ClientEndpointAddr,
				EndpointPort: holePunching.ClientEndpointPort,
				PublicKey:    publicKey[:],
				FriendlyName: *friendlyName,
				PresharedKey: preshared[:],
				Routing: &proto.RoutingInformation{
					Source:      c.identity,
					Destination: event.Routing.Source,
				},
			},
		},
	})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("peer-id", source).
			Msg("failed to send the initialization response to the remote peer")
		peer.Dispose()
		return
	}
}
