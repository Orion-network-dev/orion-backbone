package implementation

import (
	"context"
	"flag"
	"net"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"golang.zx2c4.com/wireguard/wgctrl"
	"google.golang.org/grpc"
)

var (
	friendlyName = flag.String("friendly-name", "Standard OrionD Instance", "The name of the software connecting to Orion")
	keepAlive    = flag.Duration("wireguard-keepalive", time.Second*60, "")

	allIPRanges = net.IPNet{
		IP:   net.IPv4(0, 0, 0, 0),
		Mask: net.CIDRMask(0, 32),
	}
)

// OrionClientDaemon represents the Orion client daemon with necessary components and state.
type OrionClientDaemon struct {
	// Metadata regarding the client
	memberId     uint32
	asn          uint32
	friendlyName string

	// Runtime information
	ctx context.Context

	// Structs used to manage the state of OrionD
	wireguardTunnels map[uint32]*internal.WireguardInterface
	frrManager       *FrrConfigManager
	wgClient         *wgctrl.Client

	// GRPC Clients
	registryClient     proto.RegistryClient
	holePunchingClient proto.HolePunchingServiceClient
	registryStream     proto.Registry_SubscribeToStreamClient
}

// NewOrionClientDaemon creates and initializes a new Orion client.
func NewOrionClientDaemon(ctx context.Context, clientConnection *grpc.ClientConn) (*OrionClientDaemon, error) {
	orionClient := &OrionClientDaemon{
		registryClient:     proto.NewRegistryClient(clientConnection),
		holePunchingClient: proto.NewHolePunchingServiceClient(clientConnection),
		friendlyName:       *friendlyName,
		wireguardTunnels:   make(map[uint32]*internal.WireguardInterface),
		ctx:                ctx,
	}

	var err error

	if orionClient.wgClient, err = wgctrl.New(); err != nil {
		return nil, internal.HandleError(err, "failed to create WireGuard client")
	}

	if err := orionClient.resolveIdentity(); err != nil {
		return nil, internal.HandleError(err, "failed to resolve identity")
	}

	if orionClient.frrManager, err = NewFrrConfigManager(orionClient.asn, orionClient.memberId); err != nil {
		return nil, internal.HandleError(err, "failed to initialize FRR config manager")
	}

	if err := orionClient.initializeStream(); err != nil {
		return nil, internal.HandleError(err, "failed to initialize stream")
	}

	if err := orionClient.login(); err != nil {
		return nil, internal.HandleError(err, "failed to login")
	}

	go orionClient.startListener()

	return orionClient, nil
}

func (c *OrionClientDaemon) startListener() {
	if err := c.listener(); err != nil {
		log.Error().Err(err).Msg("listener error")
	}
}

// Dispose releases resources such as WireGuard tunnels and BGP sessions.
func (c *OrionClientDaemon) Dispose() {
	log.Info().Msg("Disposing the WireGuard tunnels")
	for _, tunnel := range c.wireguardTunnels {
		tunnel.Dispose()
	}

	log.Info().Msg("Disposing the BGP sessions")
	c.frrManager.Peers = make(map[uint32]*Peer)
	if err := c.frrManager.Update(); err != nil {
		log.Error().Err(err).Msg("failed to dispose the BGP sessions")
	}
}
