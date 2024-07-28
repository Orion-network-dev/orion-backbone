package session

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (r *Session) handle_Connect(
	event *proto.RPCClientEvent_Connect,
) error {
	connect := event.Connect

	if connect.SourcePeerId == r.meta.memberId &&
		connect.DestinationPeerId != r.meta.memberId {
		log.Debug().
			Uint32("source", r.meta.memberId).
			Uint32("destination", connect.DestinationPeerId).
			Msgf("Connect Init")

		if session := r.sessionManager.GetSession(connect.DestinationPeerId); session != nil {
			session.invitations <- connect
		} else {
			log.Error().Msgf("%d is not available", connect.DestinationPeerId)
		}
	}

	return nil
}
