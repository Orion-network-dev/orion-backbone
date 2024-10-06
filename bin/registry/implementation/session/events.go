package session

import (
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

func (c *Session) HandleClientEvent(
	Event *proto.PeersToServer,
) error {
	log.Debug().Msg("handling client event")
	switch Event.Event.(type) {
	case *proto.PeersToServer_Initiate:
		return c.handle_Initiate(Event.Event.(*proto.PeersToServer_Initiate))
	case *proto.PeersToServer_InitiateAck:
		return c.handle_InitiateResponse(Event.Event.(*proto.PeersToServer_InitiateAck))
	case *proto.PeersToServer_Terminate:
		return c.handle_Terminate(Event.Event.(*proto.PeersToServer_Terminate))
	case *proto.PeersToServer_TerminateAck:
		return c.handle_TerminateACK(Event.Event.(*proto.PeersToServer_TerminateAck))
	}

	return nil
}

func (c *Session) eventListeners() {
	newClient := c.sessionManager.newClients.Listener(1)
	disposedClient := c.sessionManager.disposedClients.Listener(1)
	defer newClient.Close()
	defer disposedClient.Close()

	for {
		select {
		// Handling the events from the newClients stream
		case newClient := <-newClient.Ch():
			if c.meta != nil {
				// When a new client joins, we send a notification message
				log.Debug().
					Uint32("new-member-id", newClient.Router.MemberId).
					Uint32("session", c.meta.MemberId).
					Msgf("notifying of new client")

				c.streamSend.Broadcast(&proto.ServerToPeers{
					Event: &proto.ServerToPeers_Connected{
						Connected: newClient,
					},
				})
			}

		case disposed := <-disposedClient.Ch():
			if c.meta != nil {
				log.Debug().
					Uint32("disposed-member-id", disposed.Router.MemberId).
					Uint32("member-id", c.meta.MemberId).
					Msg("disposing client")

				c.streamSend.Broadcast(&proto.ServerToPeers{
					Event: &proto.ServerToPeers_Disconnected{
						Disconnected: disposed,
					},
				})
			}
		// Handling the events from the invitation stream
		case invitation := <-c.initiate:
			if c.meta != nil {
				if internal.IdentityFromRouter(invitation.Routing.Destination) == internal.IdentityFromRouter(c.meta) {
					log.Debug().
						Uint64("src-member-id", internal.IdentityFromRouter(invitation.Routing.Source)).
						Uint64("dst-member-id", internal.IdentityFromRouter(invitation.Routing.Destination)).
						Msg("notifying of new session invitation")

					c.streamSend.Broadcast(&proto.ServerToPeers{
						Event: &proto.ServerToPeers_Initiate{
							Initiate: invitation,
						},
					})
				} else {
					log.Error().
						Uint64("src-member-id", internal.IdentityFromRouter(invitation.Routing.Source)).
						Uint64("dst-member-id", internal.IdentityFromRouter(invitation.Routing.Destination)).
						Uint32("routine-member-id", c.meta.MemberId).
						Msg("wrong dst-id for this routine")
				}
			}
		// Handling the events from the invitation responses
		case invitation_response := <-c.initiateACK:
			if c.meta != nil {
				if internal.IdentityFromRouter(invitation_response.Routing.Destination) == internal.IdentityFromRouter(c.meta) {
					log.Debug().
						Uint64("src-member-id", internal.IdentityFromRouter(invitation_response.Routing.Source)).
						Uint64("dst-member-id", internal.IdentityFromRouter(invitation_response.Routing.Destination)).
						Msg("notifying of new invitation request")

					c.streamSend.Broadcast(&proto.ServerToPeers{
						Event: &proto.ServerToPeers_InitiateAck{
							InitiateAck: invitation_response,
						},
					})
				} else {
					log.Error().
						Uint64("src-member-id", internal.IdentityFromRouter(invitation_response.Routing.Source)).
						Uint64("dst-member-id", internal.IdentityFromRouter(invitation_response.Routing.Destination)).
						Uint32("routine-member-id", c.meta.MemberId).
						Msg("wrong dst-id for this routine")
				}
			}

		case termination := <-c.terminate:
			if c.meta != nil {
				if internal.IdentityFromRouter(termination.Routing.Destination) == internal.IdentityFromRouter(c.meta) {
					log.Debug().
						Uint64("src-member-id", internal.IdentityFromRouter(termination.Routing.Source)).
						Uint64("dst-member-id", internal.IdentityFromRouter(termination.Routing.Destination)).
						Msg("notifying of new invitation request")

					c.streamSend.Broadcast(&proto.ServerToPeers{
						Event: &proto.ServerToPeers_Terminate{
							Terminate: termination,
						},
					})
				} else {
					log.Error().
						Uint64("src-member-id", internal.IdentityFromRouter(termination.Routing.Source)).
						Uint64("dst-member-id", internal.IdentityFromRouter(termination.Routing.Destination)).
						Uint32("routine-member-id", c.meta.MemberId).
						Msg("wrong dst-id for this routine")
				}
			}

		case terminationACK := <-c.terminateACK:
			if c.meta != nil {
				if internal.IdentityFromRouter(terminationACK.Routing.Destination) == internal.IdentityFromRouter(c.meta) {
					log.Debug().
						Uint64("src-member-id", internal.IdentityFromRouter(terminationACK.Routing.Source)).
						Uint64("dst-member-id", internal.IdentityFromRouter(terminationACK.Routing.Destination)).
						Msg("notifying of new invitation request")

					c.streamSend.Broadcast(&proto.ServerToPeers{
						Event: &proto.ServerToPeers_TerminateAck{
							TerminateAck: terminationACK,
						},
					})
				} else {
					log.Error().
						Uint64("src-member-id", internal.IdentityFromRouter(terminationACK.Routing.Source)).
						Uint64("dst-member-id", internal.IdentityFromRouter(terminationACK.Routing.Destination)).
						Uint32("routine-member-id", c.meta.MemberId).
						Msg("wrong dst-id for this routine")
				}
			}

		case <-c.Context.Done():
			return
		}
	}
}
