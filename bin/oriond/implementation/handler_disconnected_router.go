package implementation

import (
	"context"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (c *OrionClientDaemon) handleDisconnectedRouter(
	ctx context.Context,
	event *proto.RouterDisconnectedEvent,
) {
	c.tunnelsLock.Lock()
	defer c.tunnelsLock.Unlock()

	peerID := internal.IdentityFromRouter(event.Router)

	peer := c.tunnels[peerID]

	if peer == nil {
		log.Error().
			Uint64("peer-id", peerID).
			Msgf("received a removed client event, but no such tunnel initialized")
		return
	}

	peer.Dispose()
	c.tunnels[peerID] = nil
}
