package session

import (
	"context"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/teivah/broadcast"
)

type Session struct {
	invitations          chan *proto.MemberConnectEvent
	invitationsResponses chan *proto.MemberConnectResponseEvent
	meta                 *SessionMeta
	sessionManager       *SessionManager

	streamSend *broadcast.Relay[*proto.RPCServerEvent]
	ctx        context.Context
	cancel     context.CancelFunc
	sID        string

	cancelCancelation chan struct{}
}

func (c *Session) Dispose() {
	// Checking if the client is auth'ed
	if c.meta != nil {
		c.cancelCancelation = make(chan struct{})
		// wait 2 minutes before ending a session
		go func() {
			log.Debug().Msg("starting to tick for session expitation")
			timer := time.NewTimer(time.Second * 20)

			select {
			case <-c.cancelCancelation:
				return
			case <-timer.C:
				c.DisposeInstant()
			}
		}()
	}
}

func (c *Session) DisposeInstant() {
	meta := c.meta
	// we should dispose the client
	c.cancel()
	c.sessionManager.disposedClients.Notify(&proto.MemberDisconnectedEvent{
		PeerId:       meta.memberId,
		FriendlyName: meta.friendlyName,
	})
	c.sessionManager.sessionIdsMap[c.sID] = nil
	c.sessionManager.sessions[c.meta.memberId] = nil
}

func (c *Session) Restore() {
	if c.meta != nil && c.cancelCancelation != nil {
		log.Info().Uint32("uid", c.meta.memberId).Msg("Session restored")
		c.cancelCancelation <- struct{}{}
	}
}

func New(
	sessionManager *SessionManager,
) (*Session, error) {
	ctx, cancel := context.WithCancel(context.Background())

	session := &Session{
		meta:                 nil,
		invitations:          make(chan *proto.MemberConnectEvent),
		invitationsResponses: make(chan *proto.MemberConnectResponseEvent),
		sessionManager:       sessionManager,
		streamSend:           broadcast.NewRelay[*proto.RPCServerEvent](),
		ctx:                  ctx,
		cancel:               cancel,
	}

	return session, nil
}

func (c *Session) Ch() <-chan *proto.RPCServerEvent {
	return c.streamSend.Listener(10).Ch()
}
