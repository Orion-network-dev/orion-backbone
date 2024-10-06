package implementation

import (
	"context"
	"net"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (c *OrionClientDaemon) handleInitiateAck(
	ctx context.Context,
	event *proto.RouterPeerToPeerInitiateACK,
) {
	c.tunnelsLock.Lock()
	defer c.tunnelsLock.Unlock()

	source := internal.IdentityFromRouter(event.Routing.Source)

	// If we get a response to a initialization request, we check if we already created the tunnel
	peerLink := c.tunnels[source]

	if peerLink == nil || peerLink.Initialized() {
		log.Error().
			Uint64("peer-id", source).
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
			Uint64("peer-id", source).
			Msg("failed to initiate a connection to the peer after a connect initialization response")
		peerLink.Dispose()
	}

	// Ends the waiting stream in the new_client handler
	c.establishedStream.Notify(source)
}
