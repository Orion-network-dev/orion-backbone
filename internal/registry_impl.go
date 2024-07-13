package internal

import (
	"crypto/x509"
	"fmt"
	"os"
	"sync"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/teivah/broadcast"
)

type OrionRegistryImplementation struct {
	newClients     *broadcast.Relay[*proto.ClientNewOnNetworkEvent]
	rootCertPool   *x509.CertPool
	clientPool     []*Client
	clientPoolLock sync.Mutex
	proto.UnimplementedRegistryServer
}

func NewOrionRegistryImplementation() (*OrionRegistryImplementation, error) {
	// Reading the root certificate
	ca, err := os.ReadFile(*AuthorityPath)
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

	return &OrionRegistryImplementation{
		newClients:     broadcast.NewRelay[*proto.ClientNewOnNetworkEvent](),
		clientPoolLock: sync.Mutex{},
		clientPool:     make([]*Client, 255),
		rootCertPool:   root,
	}, nil
}
func (r *OrionRegistryImplementation) SubscribeToStream(subscibe_event proto.Registry_SubscribeToStreamServer) error {
	var client *Client = nil

	log.Info().Msgf("new un-initialized ")
	event, err := subscibe_event.Recv()
	if err != nil {
		return err
	}

	// In case of a initialization event
	if initialize := event.GetInitialize(); initialize != nil {
		err := Authenticate(
			initialize.TimestampSigned,
			initialize.Certificate,
			initialize.Signed,
			initialize.MemberId,
			initialize.FriendlyName,
			r.rootCertPool,
		)
		if err != nil {
			log.Error().Err(err).Msg("user failed to authenticate")
			return err
		}
		// Check if this user-id is already used/connected
		if r.clientPool[initialize.MemberId] != nil {
			err := fmt.Errorf("this member_id seems to already have a running session")
			log.Debug().Err(err).Msg("this member_id already seems to be connected")
			return err
		}

		client = &Client{
			invitations:  make(chan *proto.ClientWantToConnectToClient),
			memberId:     initialize.MemberId,
			friendlyName: initialize.FriendlyName,
		}
		client.Allocate(r)
		defer client.Free(r)

		r.newClients.Broadcast(&proto.ClientNewOnNetworkEvent{
			FriendlyName: initialize.FriendlyName,
			PeerId:       initialize.MemberId,
		})
	}

	// We start a go routine to listen for global events
	go func() {
		newClients := r.newClients.Listener(1)
		for {
			select {
			case newClient := <-newClients.Ch():
				subscibe_event.Send(&proto.RPCServerEvent{
					Event: &proto.RPCServerEvent_NewClient{
						NewClient: newClient,
					},
				})
			case invitation := <-client.invitations:
				subscibe_event.Send(&proto.RPCServerEvent{
					Event: &proto.RPCServerEvent_WantsToConnect{
						WantsToConnect: invitation,
					},
				})
			case invitation_response := <-client.invitationsResponses:
				subscibe_event.Send(&proto.RPCServerEvent{
					Event: &proto.RPCServerEvent_WantsToConnectResponse{
						WantsToConnectResponse: invitation_response,
					},
				})
			case <-subscibe_event.Context().Done():
				return
			}
		}
	}()

	// Once the user is authenticated
	for {
		event, err := subscibe_event.Recv()
		if err != nil {
			return err
		}
		if connect := event.GetConnect(); connect != nil {
			if dstClient := r.clientPool[connect.PeerId]; dstClient != nil {
				dstClient.invitations <- connect
			}
		}
		if connect_response := event.GetConnectResponse(); connect_response != nil {
			if dstClient := r.clientPool[connect_response.PeerId]; dstClient != nil {
				dstClient.invitationsResponses <- connect_response
			}
		}
	}
}
