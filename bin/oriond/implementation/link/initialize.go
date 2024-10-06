package link

import (
	"flag"
	"net"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/frr"
	"github.com/MatthieuCoder/OrionV3/internal"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	keepAlive = flag.Duration("wireguard-keepalive", time.Second*60, "")

	allIPv4Ranges = net.IPNet{
		IP:   net.IPv4(0, 0, 0, 0),
		Mask: net.CIDRMask(0, 32),
	}
	allIPv6Ranges = net.IPNet{
		IP:   net.ParseIP("::"),
		Mask: net.CIDRMask(0, 128),
	}
)

func (c *PeerLink) InitializePeerConnection(
	Endpoint *net.UDPAddr,
	PublicKey wgtypes.Key,
	PresharedKey wgtypes.Key,
) error {
	c.initialized = true

	// We update our wireguard tunnel to finalize the connection request
	err := c.wireguardTunnel.SetPeers(
		c.wgClient,
		[]wgtypes.PeerConfig{
			{
				Endpoint: &net.UDPAddr{
					IP:   Endpoint.IP,
					Port: Endpoint.Port,
				},
				PresharedKey:                &PresharedKey,
				PublicKey:                   PublicKey,
				PersistentKeepaliveInterval: keepAlive,
				AllowedIPs: []net.IPNet{
					allIPv4Ranges,
					allIPv6Ranges,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	err = c.wireguardTunnel.SetAddress(c.selfIP)
	if err != nil {
		return err
	}

	// We register our peering to frr
	c.frrManager.UpdatePeer(internal.IdentityFromRouter(c.other), &frr.Peer{
		ASN:     c.other.MemberId + 64511,
		Address: c.otherIP.IP.String(),
		Weight:  200,
	})
	err = c.frrManager.Update()
	if err != nil {
		return err
	}

	// We launch our monitoring task
	go func() {
		time.Sleep(time.Second * 5)
		c.backgroundTask()
	}()

	return nil
}
