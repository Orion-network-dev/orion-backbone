package internal

import (
	"fmt"
	"net"
)

func GetAddress(self uint32) *net.IPNet {
	mask := net.CIDRMask(64, 128)
	selfAddress := net.ParseIP(fmt.Sprintf("fe80:babe::cafe:ffff:%d", self))

	return &net.IPNet{
		IP:   selfAddress,
		Mask: mask,
	}
}
