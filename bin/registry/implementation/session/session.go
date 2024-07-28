package session

import (
	"context"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

type Session struct {
	invitations          chan *proto.MemberConnectEvent
	invitationsResponses chan *proto.MemberConnectResponseEvent
	meta                 *SessionMeta
	sessionManager       *SessionManager

	streamSend chan *proto.RPCServerEvent
	ctx        context.Context
	cancel     context.CancelFunc
	sID        string

	cancelCancelation chan struct{}
}

func (c *Session) Dispose() {
	// Checking if the client is auth'ed
	if c.meta != nil {
		meta := c.meta

		// wait 2 minutes before ending a session
		go func() {
			log.Debug().Msg("starting to tick for session expitation")
			timer := time.NewTimer(time.Second * 2)

			select {
			case <-c.cancelCancelation:
				return
			case <-timer.C:
				// we should dispose the client
				c.cancel()
				c.sessionManager.disposedClients.Notify(&proto.MemberDisconnectedEvent{
					PeerId:       meta.memberId,
					FriendlyName: meta.friendlyName,
				})
				c.sessionManager.sessionIdsMap[c.sID] = nil
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
		invitations:          make(chan *proto.MemberConnectEvent),
		invitationsResponses: make(chan *proto.MemberConnectResponseEvent),
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
