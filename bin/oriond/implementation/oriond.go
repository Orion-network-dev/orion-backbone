package implementation

import (
	"context"
	"flag"
	"net"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
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

type OrionClientDaemon struct {
	// Metadata regarding the client
	memberId     uint32
	friendlyName string
	asn          uint32

	// Structs used to manage the state of OrionD
	frrManager       *FrrConfigManager
	wireguardTunnels map[uint32]*internal.WireguardInterface
	wgClient         *wgctrl.Client

	// GRPC Clients
	registryClient     proto.RegistryClient
	holePunchingClient proto.HolePunchingServiceClient
	registryStream     proto.Registry_SubscribeToStreamClient

	// Runtime information
	ctx context.Context
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
		wireguardTunnels:   make(map[uint32]*internal.WireguardInterface),
		ctx:                Context,
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
	if frrManager, err := NewFrrConfigManager(orionClient.asn, orionClient.memberId); err == nil {
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
