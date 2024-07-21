package implementation

import (
	"context"
	"flag"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/frr"
	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/link"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"golang.zx2c4.com/wireguard/wgctrl"
	"google.golang.org/grpc"
)

var (
	friendlyName = flag.String("friendly-name", "Standard OrionD Instance", "The name of the software connecting to Orion")
)

type OrionClientDaemon struct {
	// Metadata regarding the client
	memberId     uint32
	friendlyName string
	asn          uint32

	// Structs used to manage the state of OrionD
	frrManager *frr.FrrConfigManager
	wgClient   *wgctrl.Client

	// GRPC Clients
	registryClient     proto.RegistryClient
	holePunchingClient proto.HolePunchingServiceClient
	registryStream     proto.Registry_SubscribeToStreamClient

	tunnels map[uint32]*link.PeerLink

	// Runtime information
	Context   context.Context
	ctxCancel context.CancelFunc
}

// Creates and initializes a new Orion client
func NewOrionClientDaemon(
	Context context.Context,
	ClientConnection *grpc.ClientConn,
) (*OrionClientDaemon, error) {
	ctx, cancel := context.WithCancel(Context)
	orionClient := OrionClientDaemon{
		registryClient:     proto.NewRegistryClient(ClientConnection),
		holePunchingClient: proto.NewHolePunchingServiceClient(ClientConnection),
		friendlyName:       *friendlyName,
		Context:            ctx,
		ctxCancel:          cancel,
	}

	wgClient, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	orionClient.wgClient = wgClient

	// Resolve our current identity using the data from the certificates,
	// taking the overrides into acocunt
	if err := orionClient.resolveIdentity(); err != nil {
		return nil, err
	}

	// Initializing the FRR config manager, which is used to change the bgp configuration
	if frrManager, err := frr.NewFrrConfigManager(orionClient.asn, orionClient.memberId); err == nil {
		orionClient.frrManager = frrManager
	} else {
		return nil, err
	}

	// Intialize the streams to the signaling server
	if err := orionClient.initializeStream(); err != nil {
		return nil, err
	}

	// Login to the server using to initialized sessions
	if err := orionClient.login(); err != nil {
		return nil, err
	}

	// Start the listener as a background task
	go orionClient.listener()

	return &orionClient, nil
}

// Disposing interfaces and frr peers
func (c *OrionClientDaemon) Dispose() {
	log.Info().Msg("Disposing all the client")
	for _, tunnel := range c.tunnels {
		err := tunnel.Dispose()
		if err != nil {
			log.Error().
				Err(err).
				Uint32("peer-id", tunnel.RemoteID()).
				Msg("failed to dispose the client")
		}
	}
}
