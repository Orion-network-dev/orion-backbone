package internal

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

type WireguardNetLink struct {
	netlink.Link
	Id     int
	Prefix string
}

func (r WireguardNetLink) Type() string {
	return "wireguard"
}

func (r WireguardNetLink) Attrs() *netlink.LinkAttrs {
	return &netlink.LinkAttrs{
		Name: fmt.Sprintf("%s%d", r.Prefix, r.Id),
	}
}
