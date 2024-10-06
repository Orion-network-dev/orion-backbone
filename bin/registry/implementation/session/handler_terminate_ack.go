package session

import (
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (r *Session) handle_TerminateACK(
	event *proto.PeersToServer_TerminateAck,
) error {
	terminateACK := event.TerminateAck

	source := internal.IdentityFromRouter(terminateACK.Routing.Source)
	destination := internal.IdentityFromRouter(terminateACK.Routing.Destination)
	self := internal.IdentityFromRouter(r.meta)

	if source == self &&
		destination != self {
		log.Debug().
			Uint64("source", source).
			Uint64("destination", destination).
			Msgf("Connect Response")

		if session := r.sessionManager.GetSession(destination); session != nil {
			session.terminateACK <- terminateACK
		} else {
			log.Error().Msgf("%d is not available", destination)
		}
	}

	return nil
}
