package internal

import (
	"github.com/vishvananda/netlink"
)

type WireguardNetLink struct {
	netlink.Link
	InterfaceAttrs *netlink.LinkAttrs
}

func (r WireguardNetLink) Type() string {
	return "wireguard"
}

func (r WireguardNetLink) Attrs() *netlink.LinkAttrs {
	return r.InterfaceAttrs
}
