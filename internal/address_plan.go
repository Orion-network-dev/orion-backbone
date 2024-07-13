package internal

import (
	"net"

	"github.com/praserx/ipconv"
)

func szudzikPairing(x int, y int) int {
	if x < y {
		x, y = y, x
	}

	return (x * x) + x + y
}

var (
	baseIp = net.IPv4(172, 30, 0, 0)
)

func GetSelfAddress(self int, other int) (*net.IPNet, error) {
	peer := szudzikPairing(self, other)

	ipInt, err := ipconv.IPv4ToInt(baseIp)
	if err != nil {
		return nil, err
	}
	baseIPAddressInt := ipInt + uint32(peer)

	selfIPAddress := baseIPAddressInt
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
