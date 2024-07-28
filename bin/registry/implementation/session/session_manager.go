package session

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/teivah/broadcast"
)

type SessionManager struct {
	sessions map[uint32]*Session

	newClients      *broadcast.Relay[*proto.ClientNewOnNetworkEvent]
	disposedClients *broadcast.Relay[*proto.ClientDisconnectedTeardownEvent]
}

func (c *SessionManager) GetSession(session uint32) *Session {
	return c.sessions[session]
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions:        map[uint32]*Session{},
		newClients:      broadcast.NewRelay[*proto.ClientNewOnNetworkEvent](),
		disposedClients: broadcast.NewRelay[*proto.ClientDisconnectedTeardownEvent](),
	}
}
