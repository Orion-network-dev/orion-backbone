package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"google.golang.org/grpc"
)

func main() {

	// Get TLS credentials
	cred := internal.NewClientTLS()

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", "reg.orionet.re", 6443), grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalf("Unable to connect gRPC channel %v", err)
	}

	// Close the listener when containing function terminates
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("Unable to close gRPC channel %v", err)
		}
	}()

	// Create the gRPC client
	service := proto.NewDemoServiceClient(conn)

	service.SayHello(context.Background(), &proto.HelloRequest{})
}
