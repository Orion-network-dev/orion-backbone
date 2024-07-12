package main

import (
	"log"
	"net"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"golang.zx2c4.com/wireguard/wgctrl"
	"google.golang.org/grpc"
)

func main() {

	// Get TLS credentials
	cred := internal.NewServerTLS()

	// Create a listener that listens to localhost port 8443
	lis, err := net.Listen("tcp", ":6443")

	if err != nil {
		log.Fatalf("Failed to start listener %v", err)
	}

	// Close the listener when containing function terminates
	defer func() {
		err = lis.Close()
		if err != nil {
			log.Printf("Failed to close listener %v", err)
		}
	}()

	// Create a new gRPC server
	s := grpc.NewServer(grpc.Creds(cred))
	wg, err := wgctrl.New()

	proto.RegisterRegistryServer(s, &internal.OrionRegistryImplementations{})
	proto.RegisterHolePunchingServiceServer(s, &internal.OrionHolePunchingImplementations{
		WgClient: wg,
	})

	// Start the gRPC server
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
