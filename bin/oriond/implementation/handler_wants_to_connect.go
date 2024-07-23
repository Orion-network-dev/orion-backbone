package implementation

import (
	"context"
	"net"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/link"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (c *OrionClientDaemon) handleWantsToConnect(
	ctx context.Context,
	event *proto.ClientWantToConnectToClient,
) {
	// We ignore requests that are coming from ourself
	if event.DestinationPeerId != c.memberId || event.SourcePeerId == c.memberId {
		log.Error().
			Uint32("peer-id", event.SourcePeerId).
			Msg("received a message not destinated to this host")
		return
	}

	// It's our job to generate the pre-shared key information
	// according the the protocol.
	preshared, err := wgtypes.GenerateKey()
	if err != nil {
		log.Error().
			Err(err).
			Uint32("peer-id", event.SourcePeerId).
			Msg("failed to send the message to the server")
		return
	}

	// We initialize a new peer link object to hold this new link context
	peer, err := link.NewPeerLink(
		c.Context,
		c.memberId,
		event.SourcePeerId,
		c.wgClient,
		c.frrManager,
	)
	if err != nil {
		log.Error().
			Err(err).
			Uint32("peer-id", event.SourcePeerId).
			Msg("failed to initialize the peer link object")
		return
	}
	publicKey := peer.PublicKey()
	c.tunnels[event.SourcePeerId] = peer

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
			Uint32("peer-id", event.SourcePeerId).
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
			Uint32("peer-id", event.SourcePeerId).
			Msg("failed to initialize the connection to the peer")
		peer.Dispose()
		return
	}

	err = c.registryStream.Send(&proto.RPCClientEvent{
		Event: &proto.RPCClientEvent_ConnectResponse{
			ConnectResponse: &proto.ClientWantToConnectToClientResponse{
				EndpointAddr:      holePunching.ClientEndpointAddr,
				EndpointPort:      holePunching.ClientEndpointPort,
				PublicKey:         publicKey[:],
				FriendlyName:      *friendlyName,
				SourcePeerId:      c.memberId,
				DestinationPeerId: event.SourcePeerId,
				PresharedKey:      preshared[:],
			},
		},
	})
	if err != nil {
		log.Error().
			Err(err).
			Uint32("peer-id", event.SourcePeerId).
			Msg("failed to send the initialization response to the remote peer")
		peer.Dispose()
		return
	}
}
