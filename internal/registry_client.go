package internal

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
)

type Client struct {
	memberId             int64
	friendlyName         string
	invitations          chan *proto.ClientWantToConnectToClient
	invitationsResponses chan *proto.ClientWantToConnectToClientResponse
}

func (c *Client) Allocate(r *OrionRegistryImplementation) {
	r.clientPoolLock.Lock()
	defer r.clientPoolLock.Unlock()
	r.clientPool[c.memberId] = c
	log.Debug().Int64("client-id", c.memberId).Msg("Alloc client")
}

func (c *Client) Dispose(r *OrionRegistryImplementation) {
	r.clientPoolLock.Lock()
	defer r.clientPoolLock.Unlock()
	r.clientPool[c.memberId] = nil

	r.disposedClients.Broadcast(&proto.ClientDisconnectedTeardownEvent{
		PeerId:       c.memberId,
		FriendlyName: c.friendlyName,
	})

	log.Debug().Int64("client-id", c.memberId).Msg("Dealloc client")
}

func NewClient(MemberId int64, FriendlyName string) *Client {
	return &Client{
		invitations:          make(chan *proto.ClientWantToConnectToClient),
		invitationsResponses: make(chan *proto.ClientWantToConnectToClientResponse),
		memberId:             MemberId,
		friendlyName:         FriendlyName,
	}
}
