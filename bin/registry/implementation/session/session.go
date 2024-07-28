package session

import (
	"context"
	"time"

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

	cancelCancelation chan struct{}
}

func (c *Session) Dispose() {
	// Checking if the client is auth'ed
	if c.meta != nil {
		meta := c.meta

		// wait 2 minutes before ending a session
		go func() {
			timer := time.NewTimer(time.Second * 2)

			select {
			case <-c.cancelCancelation:
				return
			case <-timer.C:
				// we should dispose the client
				c.cancel()
				c.sessionManager.disposedClients.Broadcast(&proto.ClientDisconnectedTeardownEvent{
					PeerId:       meta.memberId,
					FriendlyName: meta.friendlyName,
				})
				c.sessionManager.sessions[c.meta.memberId] = nil
			}
		}()
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
