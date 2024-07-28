package session

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (c *Session) HandleClientEvent(
	Event *proto.RPCClientEvent,
) error {
	switch Event.Event.(type) {
	case *proto.RPCClientEvent_Connect:
		return c.handle_Connect(Event.Event.(*proto.RPCClientEvent_Connect))
	case *proto.RPCClientEvent_ConnectResponse:
		return c.handle_ConnectResponse(Event.Event.(*proto.RPCClientEvent_ConnectResponse))
	}
	return nil
}

func (c *Session) eventListeners() {
	newClient := c.sessionManager.newClients.Listener(1)
	disposedClient := c.sessionManager.disposedClients.Listener(1)

	for {
		select {
		// Handling the events from the newClients stream
		case newClient := <-newClient.Ch():
			// When a new client joins, we send a notification message
			log.Debug().
				Uint32("new-member-id", newClient.PeerId).
				Uint32("session", c.meta.memberId).
				Msgf("notifying of new client")

			c.streamSend <- &proto.RPCServerEvent{
				Event: &proto.RPCServerEvent_NewMember{
					NewMember: newClient,
				},
			}
		// Handling the events from the invitation stream
		case invitation := <-c.invitations:
			if invitation.DestinationPeerId == c.meta.memberId {
				log.Debug().
					Uint32("src-member-id", invitation.SourcePeerId).
					Uint32("dst-member-id", invitation.DestinationPeerId).
					Msg("notifying of new session invitation")

				c.streamSend <- &proto.RPCServerEvent{
					Event: &proto.RPCServerEvent_MemberConnect{
						MemberConnect: invitation,
					},
				}
			} else {
				log.Error().
					Uint32("src-member-id", invitation.SourcePeerId).
					Uint32("dst-member-id", invitation.DestinationPeerId).
					Uint32("routine-member-id", c.meta.memberId).
					Msg("wrong dst-id for this routine")
			}
		// Handling the events from the invitation responses
		case invitation_response := <-c.invitationsResponses:
			if invitation_response.DestinationPeerId == c.meta.memberId {
				log.Debug().
					Uint32("src-member-id", invitation_response.SourcePeerId).
					Uint32("dst-member-id", c.meta.memberId).
					Msg("notifying of new invitation request")

				c.streamSend <- &proto.RPCServerEvent{
					Event: &proto.RPCServerEvent_MemberConnectResponse{
						MemberConnectResponse: invitation_response,
					},
				}
			} else {
				log.Error().
					Uint32("src-member-id", invitation_response.SourcePeerId).
					Uint32("dst-member-id", invitation_response.DestinationPeerId).
					Uint32("routine-member-id", c.meta.memberId).
					Msg("wrong dst-id for this routine")
			}

		case disposed := <-disposedClient.Ch():
			log.Debug().
				Uint32("disposed-member-id", disposed.PeerId).
				Uint32("member-id", c.meta.memberId).
				Msg("disposing client")

			c.streamSend <- &proto.RPCServerEvent{
				Event: &proto.RPCServerEvent_DisconnectedMember{
					DisconnectedMember: disposed,
				},
			}
		case <-c.ctx.Done():
			return
		}
	}
}
