package internal

import (
	"fmt"
	"net"
)

func GetSelfAddress(self uint32, other uint32) (*net.IPNet, *net.IPNet, error) {
	mask := net.CIDRMask(64, 128)
	selfAddress := net.ParseIP(fmt.Sprintf("fe80:babe::cafe:ffff:%d", self))
	otherAddress := net.ParseIP(fmt.Sprintf("fe80:babe::cafe:ffff:%d", other))

	return &net.IPNet{
			IP:   selfAddress,
			Mask: mask,
		}, &net.IPNet{
			IP:   otherAddress,
			Mask: mask,
		}, nil
}
