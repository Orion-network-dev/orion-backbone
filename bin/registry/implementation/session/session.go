package session

import (
	"context"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/teivah/broadcast"
)

type Session struct {
	initiate     chan *proto.RouterPeerToPeerInitiate
	initiateACK  chan *proto.RouterPeerToPeerInitiateACK
	terminate    chan *proto.RouterPeerToPeerTerminate
	terminateACK chan *proto.RouterPeerToPeerTerminateACK

	meta           *proto.Router
	sessionManager *SessionManager

	streamSend *broadcast.Relay[*proto.ServerToPeers]
	Context    context.Context
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
			log.Debug().Uint32("uid", c.meta.MemberId).Msg("starting to tick for session expitation")
			timer := time.NewTimer(time.Second * 20)

			select {
			case <-c.cancelCancelation:
				c.cancelCancelation = nil
				return
			case <-timer.C:
				c.cancelCancelation = nil
				c.DisposeInstant()
			}
		}()
	}
}

func (c *Session) DisposeInstant() {
	if c.cancelCancelation != nil {
		c.cancelCancelation <- struct{}{}
	}

	meta := c.meta
	// we should dispose the client
	c.cancel()
	c.sessionManager.disposedClients.Notify(&proto.RouterDisconnectedEvent{
		Router: meta,
	})

	c.sessionManager.sessionIdsMap[c.sID] = nil
	c.sessionManager.sessions[internal.IdentityFromRouter(c.meta)] = nil

	// todo: implements ack in the protocol
	time.Sleep(2 * time.Second)
}

func (c *Session) Restore() {
	if c.meta != nil && c.cancelCancelation != nil {
		log.Info().Uint32("uid", c.meta.MemberId).Msg("Session restored")
		c.cancelCancelation <- struct{}{}
	}
}

func New(
	sessionManager *SessionManager,
) *Session {
	ctx, cancel := context.WithCancel(context.Background())

	session := &Session{
		meta:           nil,
		initiate:       make(chan *proto.RouterPeerToPeerInitiate),
		initiateACK:    make(chan *proto.RouterPeerToPeerInitiateACK),
		terminate:      make(chan *proto.RouterPeerToPeerTerminate),
		terminateACK:   make(chan *proto.RouterPeerToPeerTerminateACK),
		sessionManager: sessionManager,
		streamSend:     broadcast.NewRelay[*proto.ServerToPeers](),
		Context:        ctx,
		cancel:         cancel,
	}

	return session
}

func (c *Session) Ch() *broadcast.Listener[*proto.ServerToPeers] {
	return c.streamSend.Listener(10)
}
