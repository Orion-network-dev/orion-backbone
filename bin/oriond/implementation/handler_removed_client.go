package implementation

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (c *OrionClientDaemon) handleRemovedClient(event *proto.ClientDisconnectedTeardownEvent) {
	wg := c.wireguardTunnels[event.PeerId]

	// If we had a wireguard tunnel to this peer, we shall delete it
	if wg != nil {
		wg.Dispose()
	}

	if peer := c.frrManager.Peers[event.PeerId]; peer != nil {
		c.frrManager.Peers[event.PeerId] = nil
		if err := c.frrManager.Update(); err != nil {
			log.Error().Err(err).Msg("failed to apply frr configuration changes")
		}
	}
}
