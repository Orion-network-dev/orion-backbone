package implementation

import (
	"context"
	"flag"
	"sync"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/frr"
	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/link"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/teivah/broadcast"
	"golang.zx2c4.com/wireguard/wgctrl"
	"google.golang.org/grpc"
)

var (
	friendlyName = flag.String("friendly-name", "Standard OrionD Instance", "The name of the software connecting to Orion")
)

type OrionClientDaemon struct {
	// Metadata regarding the client
	identity     *proto.Router
	friendlyName string
	asn          uint32
	sID          string

	// Structs used to manage the state of OrionD
	frrManager *frr.FrrConfigManager
	wgClient   *wgctrl.Client

	// GRPC Clients
	registryClient     proto.RegistryClient
	holePunchingClient proto.HolePunchingServiceClient
	registryStream     proto.Registry_SubscribeToStreamClient

	tunnels     map[uint64]*link.PeerLink
	tunnelsLock *sync.RWMutex

	// Runtime information
	Context           context.Context
	ParentCtx         context.Context
	ctxCancel         context.CancelFunc
	establishedStream *broadcast.Relay[uint64]
}

// Creates and initializes a new Orion client
func NewOrionClientDaemon(
	Context context.Context,
	ClientConnection *grpc.ClientConn,
) (*OrionClientDaemon, error) {
	orionClient := OrionClientDaemon{
		registryClient:     proto.NewRegistryClient(ClientConnection),
		holePunchingClient: proto.NewHolePunchingServiceClient(ClientConnection),
		friendlyName:       *friendlyName,
		ParentCtx:          Context,
		establishedStream:  broadcast.NewRelay[uint64](),
		tunnels:            map[uint64]*link.PeerLink{},
		tunnelsLock:        &sync.RWMutex{},
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
	if frrManager, err := frr.NewFrrConfigManager(orionClient.asn, orionClient.identity); err == nil {
		orionClient.frrManager = frrManager
	} else {
		return nil, err
	}

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
				Uint64("peer-id", internal.IdentityFromRouter(tunnel.RemoteID())).
				Msg("failed to dispose the client")
		}
	}
}
