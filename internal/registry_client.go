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
func (c *Client) Free(r *OrionRegistryImplementation) {
	r.clientPoolLock.Lock()
	defer r.clientPoolLock.Unlock()
	r.clientPool[c.memberId] = nil

	log.Debug().Int64("client-id", c.memberId).Msg("Dealloc client")
}
