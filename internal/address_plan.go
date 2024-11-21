package internal

import (
	"fmt"
	"net"
	"slices"
)

func GetSelfAddress(self uint32, other uint32) (*net.IPNet, *net.IPNet, error) {
	mask := net.CIDRMask(120, 128)
	nid := []uint32{self, other}
	slices.Sort(nid)

	selfAddress := net.ParseIP(fmt.Sprintf("fe80:babe::cafe::%d:%d:%d", nid[0], nid[1], self))
	otherAddress := net.ParseIP(fmt.Sprintf("fe80:babe::cafe::%d:%d:%d", nid[0], nid[1], other))

	return &net.IPNet{
			IP:   selfAddress,
			Mask: mask,
		}, &net.IPNet{
			IP:   otherAddress,
			Mask: mask,
		}, nil
}
