package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/grpc"
)

func main() {

	// Get TLS credentials
	cred := internal.NewClientTLS()

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", "reg.orionet.re", 6443), grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalf("Unable to connect gRPC channel %v", err)
	}

	// Create the gRPC client
	_ = proto.NewRegistryClient(conn)
	holepunch := proto.NewHolePunchingServiceClient(conn)

	pk, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	events, err := holepunch.Session(ctx, &proto.HolePunchingInitialize{
		PublicKey: pk.PublicKey().String(),
	})
	if err != nil {
		panic(err)
	}

	for {
		event, err := events.Recv()

		if err != nil {
			panic(err)
			return
		}
		fmt.Println(event)
	}
}
