package implementation

import (
	"net"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (c *OrionClientDaemon) handleWantsToConnectResponse(event *proto.ClientWantToConnectToClientResponse) {
	// If we get a response to a initialization request, we check if we already created the tunnel
	peerLink := c.tunnels[event.SourcePeerId]

	if peerLink == nil || peerLink.Initialized() {
		log.Error().
			Uint32("peer-id", event.SourcePeerId).
			Msg("received an invalid connect response")
		return
	}

	// If everything is allright we initialize the connection to the peer
	err := peerLink.InitializePeerConnection(
		&net.UDPAddr{
			IP:   net.ParseIP(event.EndpointAddr),
			Port: int(event.EndpointPort),
		},
		wgtypes.Key(event.PublicKey),
		wgtypes.Key(event.PresharedKey),
	)
	if err != nil {
		log.Error().
			Err(err).
			Uint32("peer-id", event.SourcePeerId).
			Msg("failed to initiate a connection to the peer after a connect initialization response")
		peerLink.Dispose()
	}

	// Ends the waiting stream in the new_client handler
	c.establishedStream <- event.SourcePeerId
}
