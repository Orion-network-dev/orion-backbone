package internal

import (
	"math/big"
	"net"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/praserx/ipconv"
	"github.com/rs/zerolog/log"
)

func szudzikPairing(x uint64, y uint64) uint64 {
	if x < y {
		x, y = y, x
	}

	return (x * x) + x + y
}

var (
	// This is the ipv6 network used for interconnects: fd0b:0b5f:2486::/48
	baseIp = net.ParseIP("fd0b:0b5f:2486::")
)

func GetSelfAddress(self *proto.Router, other *proto.Router) (*net.IPNet, *net.IPNet, error) {
	selfID := IdentityFromRouter(self)
	otherID := IdentityFromRouter(other)

	peer := szudzikPairing(selfID, otherID)

	ipInt, err := ipconv.IPv6ToBigInt(baseIp)
	if err != nil {
		log.Error().Err(err).Msgf("failed to convert to ip address to a big int interger")
		return nil, nil, err
	}

	selfIPAddress := big.NewInt(0).Add(ipInt, big.NewInt(int64(peer<<1)))
	otherIPAddress := big.NewInt(0).Add(ipInt, big.NewInt(int64(peer<<1)))

	if selfID > otherID {
		selfIPAddress = big.NewInt(0).Add(selfIPAddress, big.NewInt(1))
	} else {
		otherIPAddress = big.NewInt(0).Add(otherIPAddress, big.NewInt(1))
	}

	mask := net.CIDRMask(127, 128)
	selfIPNet := net.IPNet{
		IP:   ipconv.BigIntToIPv6(*selfIPAddress),
		Mask: mask,
	}
	otherIPNet := net.IPNet{
		IP:   ipconv.BigIntToIPv6(*otherIPAddress),
		Mask: mask,
	}

	return &selfIPNet, &otherIPNet, nil
}
