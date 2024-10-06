package link

import (
	"context"
	"fmt"
	"net"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/frr"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type PeerLink struct {
	ctx             context.Context
	frrManager      *frr.FrrConfigManager
	wireguardTunnel *internal.WireguardInterface
	wgClient        *wgctrl.Client

	publicKey   wgtypes.Key
	self        *proto.Router
	other       *proto.Router
	otherIP     *net.IPNet
	selfIP      *net.IPNet
	cancel      context.CancelFunc
	initialized bool
}

func NewPeerLink(
	parentCtx context.Context,
	self *proto.Router,
	other *proto.Router,
	wgClient *wgctrl.Client,
	frrManager *frr.FrrConfigManager,
) (*PeerLink, error) {
	selfIP, otherIP, err := internal.GetSelfAddress(self, other)
	if err != nil {
		log.Error().Err(err).Msg("failed to compute the self address")
		return nil, err
	}
	privatekey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	tunnel, err := internal.NewWireguardInterface(wgClient, &netlink.LinkAttrs{
		Name:  fmt.Sprintf("orion%d%d", other.MemberId, other.RouterId),
		Group: 30,
	}, wgtypes.Config{
		PrivateKey:   &privatekey,
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{},
	})
	if err != nil {
		return nil, err
	}
	err = tunnel.SetAddress(selfIP)
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
		self:            self,
		other:           other,
		selfIP:          selfIP,
		otherIP:         otherIP,
		publicKey:       privatekey.PublicKey(),
		cancel:          cancel,
		initialized:     false,
	}

	return link, nil
}

func (c *PeerLink) PublicKey() wgtypes.Key {
	return c.publicKey
}

func (c *PeerLink) Initialized() bool {
	return c.initialized
}

func (c *PeerLink) RemoteID() *proto.Router {
	return c.other
}
