package implementation

import (
	"context"
	"crypto/x509"
	"flag"
	"sync"

	"github.com/orion-network-dev/orion-backbone/bin/oriond/implementation/frr"
	"github.com/orion-network-dev/orion-backbone/bin/oriond/implementation/link"
	"github.com/rs/zerolog/log"
	"github.com/teivah/broadcast"
	"golang.zx2c4.com/wireguard/wgctrl"
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

	tunnels     map[uint32]*link.PeerLink
	tunnelsLock *sync.RWMutex

	// Runtime information
	Context           context.Context
	ParentCtx         context.Context
	establishedStream *broadcast.Relay[uint32]

	chain []*x509.Certificate
}

// Creates and initializes a new Orion client
func NewOrionClientDaemon(
	parentContext context.Context,
) (*OrionClientDaemon, error) {
	orionClient := OrionClientDaemon{
		friendlyName:      *friendlyName,
		ParentCtx:         parentContext,
		establishedStream: broadcast.NewRelay[uint32](),
		tunnels:           map[uint32]*link.PeerLink{},
		tunnelsLock:       &sync.RWMutex{},
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
