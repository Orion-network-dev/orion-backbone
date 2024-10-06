package internal

import (
	"github.com/MatthieuCoder/OrionV3/internal/proto"
)

func IdentityFromRouter(router *proto.Router) uint64 {
	return uint64(router.RouterId) + uint64(router.MemberId<<32)
}
