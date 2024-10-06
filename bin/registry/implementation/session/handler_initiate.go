package session

import (
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (r *Session) handle_Initiate(
	event *proto.PeersToServer_Initiate,
) error {
	initiate := event.Initiate

	source := internal.IdentityFromRouter(initiate.Routing.Source)
	destination := internal.IdentityFromRouter(initiate.Routing.Destination)
	self := internal.IdentityFromRouter(r.meta)

	if source == self &&
		destination != self {
		log.Debug().
			Uint64("source", source).
			Uint64("destination", destination).
			Msgf("Connect Init")

		if session := r.sessionManager.GetSession(destination); session != nil {
			session.initiate <- initiate
		} else {
			log.Error().Msgf("%d is not available", destination)
		}
	}

	return nil
}
