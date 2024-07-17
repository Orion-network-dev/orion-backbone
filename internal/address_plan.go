package internal

import (
	"net"

	"github.com/praserx/ipconv"
	"github.com/rs/zerolog/log"
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

func GetSelfAddress(self uint32, other uint32) (*net.IPNet, *net.IPNet, error) {
	peer := szudzikPairing(self, other)

	ipInt, err := ipconv.IPv4ToInt(baseIp)
	if err != nil {
		log.Error().Err(err).Msgf("failed to convert to ip address to a uint32 interger")
		return nil, nil, err
	}

	selfIPAddress := ipInt + uint32(peer<<1)
	otherIPAddress := ipInt + uint32(peer<<1)

	if self > other {
		selfIPAddress = selfIPAddress + 1
	} else {
		otherIPAddress = otherIPAddress + 1
	}

	mask := net.CIDRMask(31, 32)
	selfIPNet := net.IPNet{
		IP:   ipconv.IntToIPv4(selfIPAddress),
		Mask: mask,
	}
	otherIPNet := net.IPNet{
		IP:   ipconv.IntToIPv4(otherIPAddress),
		Mask: mask,
	}

	return &selfIPNet, &otherIPNet, nil
}
