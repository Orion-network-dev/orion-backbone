package internal

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"
	"sync"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/teivah/broadcast"
)

type Client struct {
	memberId     int64
	friendlyName string
	invitations  chan *proto.ClientWantsToInitiateLinkEvent
}

func (c *Client) Allocate(r *OrionRegistryImplementations) {
	r.clientPoolLock.Lock()
	defer r.clientPoolLock.Unlock()
	r.clientPool[c.memberId] = c
}
func (c *Client) Free(r *OrionRegistryImplementations) {
	r.clientPoolLock.Lock()
	defer r.clientPoolLock.Unlock()
	r.clientPool[c.memberId] = nil
}

type OrionRegistryImplementations struct {
	newClients     *broadcast.Relay[*proto.ClientNewOnNetworkEvent]
	rootCertPool   *x509.CertPool
	clientPool     []*Client
	clientPoolLock sync.Mutex
	proto.UnimplementedRegistryServer
}

func NewOrionRegistryImplementation() (*OrionRegistryImplementations, error) {
	// Reading the root certificate
	ca, err := os.ReadFile("ca/ca.crt")
	if err != nil {
		log.Debug().Err(err).Msg("failed to import the root ca certificate")
		return nil, err
	}

	// Create a new certificate pool containing the root certificates
	root := x509.NewCertPool()
	// Append the root certificate to the pool
	ok := root.AppendCertsFromPEM(ca)
	if !ok {
		return nil, fmt.Errorf("the root certificate failed to be imported")
	}

	return &OrionRegistryImplementations{
		newClients:     broadcast.NewRelay[*proto.ClientNewOnNetworkEvent](),
		clientPoolLock: sync.Mutex{},
		clientPool:     make([]*Client, 255),
		rootCertPool:   root,
	}, nil
}

func (r *OrionRegistryImplementations) SubscribeToStream(initializeRequest *proto.InitializeRequest, z proto.Registry_SubscribeToStreamServer) error {
	// Call the Authentication function to check all the user-given parameters gainst the signature and root-ca
	err := Authenticate(
		initializeRequest.TimestampSigned,
		initializeRequest.Certificate,
		initializeRequest.Signed,
		initializeRequest.MemberId,
		initializeRequest.FriendlyName,
		r.rootCertPool,
	)
	if err != nil {
		log.Error().Err(err).Msg("user failed to authenticate")
		return err
	}

	// Check if this user-id is already used/connected
	if r.clientPool[initializeRequest.MemberId] != nil {
		err := fmt.Errorf("this member_id seems to already have a running session")
		log.Debug().Err(err).Msg("this member_id already seems to be connected")
		return err
	}

	// Since this user-id is free, we process to allocate it.
	client := &Client{
		invitations:  make(chan *proto.ClientWantsToInitiateLinkEvent),
		memberId:     initializeRequest.MemberId,
		friendlyName: initializeRequest.FriendlyName,
	}
	client.Allocate(r)
	defer client.Free(r)

	// Tell the other Orion network members that a new client has arrived
	r.newClients.Broadcast(&proto.ClientNewOnNetworkEvent{
		FriendlyName: initializeRequest.FriendlyName,
		PeerId:       initializeRequest.MemberId,
	})

	// Subscribe to the new clients stream.
	listener := r.newClients.Listener(1)
	newClientsEvents := listener.Ch()

	for {
		select {

		case newClientData := <-newClientsEvents:
			if newClient := r.clientPool[newClientData.PeerId]; newClient != nil {
				newClient.invitations <- &proto.ClientWantsToInitiateLinkEvent{
					FriendlyName: client.friendlyName,
					PeerId:       client.memberId,
				}
			}

		case initiateRequest := <-client.invitations:
			z.Send(&proto.RPCEvent{
				Event: &proto.RPCEvent_ClientWantsToInitiateLinkEvent{
					ClientWantsToInitiateLinkEvent: initiateRequest,
				},
			})
		case <-z.Context().Done():
			return z.Context().Err()
		}
	}
}

// When an existing client wants to initiate a connection to a new or existing peer.
func (r *OrionHolePunchingImplementation) InitializeConnectionToPeer(context.Context, *proto.InitializeConnectionToPeerRequest) (*proto.InitializeConnectionToPeerResponse, error) {
	return nil, nil
}
