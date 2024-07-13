package internal

import "github.com/MatthieuCoder/OrionV3/internal/proto"

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
}
func (c *Client) Free(r *OrionRegistryImplementation) {
	r.clientPoolLock.Lock()
	defer r.clientPoolLock.Unlock()
	r.clientPool[c.memberId] = nil
}
