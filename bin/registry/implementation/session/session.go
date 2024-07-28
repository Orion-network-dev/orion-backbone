package session

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
)

type Session struct {
	invitations          chan *proto.ClientWantToConnectToClient
	invitationsResponses chan *proto.ClientWantToConnectToClientResponse
	meta                 *SessionMeta
	sessionManager       *SessionManager

	streamSend chan *proto.RPCServerEvent
}

func (c *Session) Dispose() {
	// Checking if the client is auth'ed
	if c.meta != nil {
		meta := c.meta

		// Inform all other clients that a client is dead
		c.sessionManager.disposedClients.Broadcast(&proto.ClientDisconnectedTeardownEvent{
			PeerId:       meta.memberId,
			FriendlyName: meta.friendlyName,
		})
		c.sessionManager.sessions[c.meta.memberId] = nil
	}
}

func New(
	sessionManager *SessionManager,
) (*Session, error) {

	session := &Session{
		meta:                 nil,
		invitations:          make(chan *proto.ClientWantToConnectToClient),
		invitationsResponses: make(chan *proto.ClientWantToConnectToClientResponse),
		sessionManager:       sessionManager,
		streamSend:           make(chan *proto.RPCServerEvent),
	}

	return session, nil
}

func (c *Session) Ch() chan *proto.RPCServerEvent {
	return c.streamSend
}
