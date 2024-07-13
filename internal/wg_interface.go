package internal

import (
	"github.com/rs/zerolog/log"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type WireguardInterface struct {
	wglink   WireguardNetLink
	wgconfig wgtypes.Config
}

func NewWireguardInterface(wg *wgctrl.Client, interfaceAttrs *netlink.LinkAttrs, configuration wgtypes.Config) (*WireguardInterface, error) {

	wglink := WireguardNetLink{
		InterfaceAttrs: interfaceAttrs,
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

	return &WireguardInterface{
		wglink: wglink,
	}, nil
}

func (c *WireguardInterface) SetPeers(wg *wgctrl.Client, peers []wgtypes.PeerConfig) error {
	c.wgconfig.Peers = peers
	c.wgconfig.ReplacePeers = true

	if err := wg.ConfigureDevice(c.wglink.InterfaceAttrs.Name, c.wgconfig); err != nil {
		log.Error().Err(err).Msg("failed to apply the wireguard configuration")
		netlink.LinkDel(c.wglink)
		return err
	}
	return nil
}

func (c *WireguardInterface) Dispose() {
	if err := netlink.LinkDel(c.wglink); err != nil {
		log.Error().Err(err).Msg("failed to destroy interface")
	}
}
