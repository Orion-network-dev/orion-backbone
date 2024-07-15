package internal

import (
	"net"

	"github.com/praserx/ipconv"
)

func szudzikPairing(x uint32, y uint32) uint32 {
	if x < y {
		x, y = y, x
	}

	return (x * x) + x + y
}

var (
	baseIp = net.IPv4(172, 30, 0, 0)
)

func GetSelfAddress(self uint32, other uint32) (*net.IPNet, error) {
	peer := szudzikPairing(self, other)

	ipInt, err := ipconv.IPv4ToInt(baseIp)
	if err != nil {
		return nil, err
	}
	selfIPAddress := ipInt + uint32(peer<<1)

	if self > other {
		selfIPAddress = selfIPAddress + 1
	}

	mask := net.CIDRMask(31, 32)
	selfIPNet := net.IPNet{
		IP:   ipconv.IntToIPv4(selfIPAddress),
		Mask: mask,
	}

	return &selfIPNet, nil
}
