package session

import (
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (r *Session) handle_Terminate(
	event *proto.PeersToServer_Terminate,
) error {
	terminate := event.Terminate

	source := internal.IdentityFromRouter(terminate.Routing.Source)
	destination := internal.IdentityFromRouter(terminate.Routing.Destination)
	self := internal.IdentityFromRouter(r.meta)

	if source == self &&
		destination == self {
		log.Debug().
			Uint64("source", source).
			Uint64("destination", destination).
			Msgf("Connect Response")

		if session := r.sessionManager.GetSession(destination); session != nil {
			session.terminate <- terminate
		} else {
			log.Error().Msgf("%d is not available", destination)
		}
	}

	return nil
}
