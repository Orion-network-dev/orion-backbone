package link

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/orion-network-dev/orion-backbone/bin/oriond/implementation/frr"
	"github.com/orion-network-dev/orion-backbone/internal"
	"github.com/vishvananda/netlink"
	"gitlab.com/NebulousLabs/go-upnp"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	basePort = flag.Uint("override-base-port", 65000, "Override the public port for this instance")
)

type PeerLink struct {
	ctx             context.Context
	frrManager      *frr.FrrConfigManager
	wireguardTunnel *internal.WireguardInterface
	wgClient        *wgctrl.Client

	publicKey   wgtypes.Key
	selfIP      *net.IPNet
	otherIP     *net.IPNet
	selfID      uint32
	otherID     uint32
	cancel      context.CancelFunc
	initialized bool
	externalIP  *string
	igd         *upnp.IGD
}

func NewPeerLink(
	parentCtx context.Context,
	selfID uint32,
	otherID uint32,
	wgClient *wgctrl.Client,
	frrManager *frr.FrrConfigManager,
) (*PeerLink, error) {
	// we get our link-local address according to the adressing plan
	ipAddress := internal.GetAddress(selfID)

	// we generate the unique tunnel wireguard privateKey
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	port := int(*basePort + uint(otherID))
	tunnel, err := internal.NewWireguardInterface(wgClient, &netlink.LinkAttrs{
		Name:  fmt.Sprintf("orion%d", otherID),
		Group: 30,
	}, wgtypes.Config{
		PrivateKey:   &privateKey,
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{},
		ListenPort:   &port,
	})
	if err != nil {
		return nil, err
	}

	err = tunnel.SetAddress(ipAddress)
	if err != nil {
		defer tunnel.Dispose()
		return nil, err
	}

	ctx, cancel := context.WithCancel(parentCtx)
	link := &PeerLink{
		ctx:             ctx,
		frrManager:      frrManager,
		wireguardTunnel: tunnel,
		wgClient:        wgClient,
		selfID:          selfID,
		otherID:         otherID,
		selfIP:          ipAddress,
		publicKey:       privateKey.PublicKey(),
		cancel:          cancel,
		initialized:     false,
	}
	link.upnpInit()

	return link, nil
}

func (c *PeerLink) PublicKey() wgtypes.Key {
	return c.publicKey
}

func (c *PeerLink) Initialized() bool {
	return c.initialized
}

func (c *PeerLink) RemoteID() uint32 {
	return c.otherID
}
