package implementation

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (c *OrionClientDaemon) handleRemovedClient(event *proto.ClientDisconnectedTeardownEvent) {
	peer := c.tunnels[event.PeerId]

	if peer == nil {
		log.Error().
			Uint32("peer-id", event.PeerId).
			Msgf("received a removed client event, but no such tunnel initialized")
		return
	}
	// Since every ressource link to a peer is linked to a PeerLink
	// we simply have to dispose the peer to remove all resources
	peer.Dispose()
}
