package link

import (
	"flag"
	"net"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/frr"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	keepAlive = flag.Duration("wireguard-keepalive", time.Second*60, "")

	allIPRanges = net.IPNet{
		IP:   net.IPv4(0, 0, 0, 0),
		Mask: net.CIDRMask(0, 32),
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
					allIPRanges,
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
	c.frrManager.UpdatePeer(c.otherID, &frr.Peer{
		ASN:     c.otherID + 64511,
		Address: c.otherIP.IP.String(),
		OrionId: c.otherID,
	})
	c.wireguardTunnel.AddRoute(c.otherID, 1)
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
