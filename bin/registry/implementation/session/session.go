package session

import (
	"context"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
)

type Session struct {
	invitations          chan *proto.ClientWantToConnectToClient
	invitationsResponses chan *proto.ClientWantToConnectToClientResponse
	meta                 *SessionMeta
	sessionManager       *SessionManager

	streamSend chan *proto.RPCServerEvent
	ctx        context.Context
	cancel     context.CancelFunc
}

func (c *Session) Dispose() {
	// Checking if the client is auth'ed
	if c.meta != nil {
		meta := c.meta
		c.cancel()

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
	ctx, cancel := context.WithCancel(context.Background())

	session := &Session{
		meta:                 nil,
		invitations:          make(chan *proto.ClientWantToConnectToClient),
		invitationsResponses: make(chan *proto.ClientWantToConnectToClientResponse),
		sessionManager:       sessionManager,
		streamSend:           make(chan *proto.RPCServerEvent),
		ctx:                  ctx,
		cancel:               cancel,
	}

	return session, nil
}

func (c *Session) Ch() chan *proto.RPCServerEvent {
	return c.streamSend
}
