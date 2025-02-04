package internal

import (
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type WireguardInterface struct {
	WgLink   WireguardNetLink
	WgConfig wgtypes.Config
	lock     sync.Mutex
}

func NewWireguardInterface(wg *wgctrl.Client, interfaceAttrs *netlink.LinkAttrs, configuration wgtypes.Config) (*WireguardInterface, error) {
	log.Debug().Str("interface", interfaceAttrs.Name).Msg("configuring interface")
	wglink := WireguardNetLink{
		InterfaceAttrs: interfaceAttrs,
	}

	if link, err := netlink.LinkByName(wglink.InterfaceAttrs.Name); err == nil {
		log.Debug().Str("interface", interfaceAttrs.Name).Msg("interface already exist, deleting")
		if err := netlink.LinkDel(link); err != nil {
			log.Error().Msg("failed to delete the already existing interface")
		}
	}

	if err := netlink.LinkAdd(wglink); err != nil {
		log.Error().Err(err).Msg("error while creating the interface")
		return nil, err
	}

	if err := wg.ConfigureDevice(wglink.InterfaceAttrs.Name, configuration); err != nil {
		log.Error().Err(err).Msg("failed to apply the wireguard configuration")
		netlink.LinkDel(wglink)
		return nil, err
	}

	if err := netlink.LinkSetUp(wglink); err != nil {
		log.Error().Err(err).Msg("failed to set the interface up")
		netlink.LinkDel(wglink)
		return nil, err
	}

	log.Debug().Str("interface", interfaceAttrs.Name).Msg("finished setting up interface")
	return &WireguardInterface{
		WgLink: wglink,
	}, nil
}

func (c *WireguardInterface) SetPeers(wg *wgctrl.Client, peers []wgtypes.PeerConfig) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debug().Str("interface", c.WgLink.InterfaceAttrs.Name).Msg("updating peers on interface")
	c.WgConfig.Peers = peers
	c.WgConfig.ReplacePeers = true

	if err := wg.ConfigureDevice(c.WgLink.InterfaceAttrs.Name, c.WgConfig); err != nil {
		log.Error().Err(err).Msg("failed to apply the wireguard configuration")
		netlink.LinkDel(c.WgLink)
		return err
	}
	return nil
}

func (c *WireguardInterface) Dispose() {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debug().Str("interface", c.WgLink.InterfaceAttrs.Name).Msg("disposing wireguard interface")
	if err := netlink.LinkDel(c.WgLink); err != nil {
		log.Error().Err(err).Msg("failed to destroy interface")
	}
}

func (c *WireguardInterface) SetAddress(ip *net.IPNet) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debug().Str("interface", c.WgLink.InterfaceAttrs.Name).Msg("updating the IP address")

	existingIPs, err := netlink.AddrList(c.WgLink, netlink.FAMILY_ALL)
	if err != nil {
		return err
	}

	// Check if we already have the address
	for _, existingIP := range existingIPs {
		if existingIP.IP.Equal(ip.IP) {
			return nil
		}
	}

	// Otherwise we add the address
	if err := netlink.AddrAdd(c.WgLink, &netlink.Addr{
		IPNet: ip,
	}); err != nil {
		log.Error().Err(err).Msg("failed to assign IP addresses")
		return err
	}

	return nil
}
