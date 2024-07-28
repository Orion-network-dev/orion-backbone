package session

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/teivah/broadcast"
)

type SessionManager struct {
	sessions      map[uint32]*Session
	sessionIdsMap map[string]*uint32

	newClients      *broadcast.Relay[*proto.NewMemberEvent]
	disposedClients *broadcast.Relay[*proto.MemberDisconnectedEvent]
}

func (c *SessionManager) GetSession(session uint32) *Session {
	return c.sessions[session]
}

func (c *SessionManager) GetSessionFromSessionId(id string) *Session {
	uid := c.sessionIdsMap[id]
	if uid == nil {
		return nil
	}
	return c.GetSession(*uid)
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions:        map[uint32]*Session{},
		newClients:      broadcast.NewRelay[*proto.NewMemberEvent](),
		disposedClients: broadcast.NewRelay[*proto.MemberDisconnectedEvent](),
	}
}
