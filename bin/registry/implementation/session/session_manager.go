package session

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/teivah/broadcast"
)

type SessionManager struct {
	sessions      map[uint64]*Session
	sessionIdsMap map[string]*uint64

	newClients      *broadcast.Relay[*proto.RouterConnectedEvent]
	disposedClients *broadcast.Relay[*proto.RouterDisconnectedEvent]
}

func (c *SessionManager) GetSession(session uint64) *Session {
	return c.sessions[session]
}

func (c *SessionManager) GetSessionFromSessionId(id string) *Session {
	log.Debug().Str("session-id", id).Msg("getting by session id")
	uid := c.sessionIdsMap[id]
	if uid == nil {
		return nil
	}
	return c.GetSession(*uid)
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions:        map[uint64]*Session{},
		sessionIdsMap:   make(map[string]*uint64),
		newClients:      broadcast.NewRelay[*proto.RouterConnectedEvent](),
		disposedClients: broadcast.NewRelay[*proto.RouterDisconnectedEvent](),
	}
}
