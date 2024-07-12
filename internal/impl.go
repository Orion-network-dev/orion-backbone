package internal

import (
	"context"
	"fmt"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
)

type OrionRegistryImpl struct {
	proto.DemoServiceServer
}

func (r *OrionRegistryImpl) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloResponse, error) {
	fmt.Println("Hello given!!")
	return &proto.HelloResponse{}, nil
}
