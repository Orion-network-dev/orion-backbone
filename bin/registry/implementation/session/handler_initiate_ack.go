package session

import (
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (r *Session) handle_InitiateResponse(
	event *proto.PeersToServer_InitiateAck,
) error {
	initiateACK := event.InitiateAck

	source := internal.IdentityFromRouter(initiateACK.Routing.Source)
	destination := internal.IdentityFromRouter(initiateACK.Routing.Destination)
	self := internal.IdentityFromRouter(r.meta)

	if source == self &&
		destination != self {
		log.Debug().
			Uint64("source", source).
			Uint64("destination", destination).
			Msgf("Connect Response")

		if session := r.sessionManager.GetSession(destination); session != nil {
			session.initiateACK <- initiateACK
		} else {
			log.Error().Msgf("%d is not available", destination)
		}
	}

	return nil
}
