package session

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (r *Session) handle_ConnectResponse(
	event *proto.RPCClientEvent_ConnectResponse,
) error {
	connectResponse := event.ConnectResponse

	if connectResponse.SourcePeerId == r.meta.memberId &&
		connectResponse.DestinationPeerId != r.meta.memberId {
		log.Debug().
			Uint32("source", r.meta.memberId).
			Uint32("destination", connectResponse.DestinationPeerId).
			Msgf("Connect Response")

		if session := r.sessionManager.GetSession(connectResponse.DestinationPeerId); session != nil {
			session.invitationsResponses <- connectResponse
		} else {
			log.Error().Msgf("%d is not available", connectResponse.DestinationPeerId)
		}
	}

	return nil
}
