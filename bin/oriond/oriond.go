package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"encoding/pem"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"google.golang.org/grpc"
)

func main() {

	// Get TLS credentials
	cred := internal.NewClientTLS()

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", "reg.orionet.re", 6443), grpc.WithTransportCredentials(cred), grpc.WithIdleTimeout(time.Second*120))
	if err != nil {
		log.Fatalf("Unable to connect gRPC channel %v", err)
	}

	// Create the gRPC client
	registryClient := proto.NewRegistryClient(conn)
	_ = proto.NewHolePunchingServiceClient(conn)
	now := time.Now()
	sec := now.Unix()
	digest := fmt.Sprintf("%d:%s:%d", 0, "0.orionet.re", sec)
	hash := sha512.New().Sum([]byte(digest))
	certPEM, err := os.ReadFile("ca/client.crt")
	if err != nil {
		panic(err)
	}

	privateKey, err := os.ReadFile("ca/client.key")
	if err != nil {
		panic(err)
	}
	zzz, _ := pem.Decode(privateKey)

	pk, err := x509.ParseECPrivateKey(zzz.Bytes)
	if err != nil {
		panic(err)
	}
	signature, err := ecdsa.SignASN1(rand.Reader, pk, hash)
	if err != nil {
		panic(err)
	}
	data := &proto.InitializeRequest{
		FriendlyName:    "0.orionet.re",
		TimestampSigned: sec,
		MemberId:        0,
		Certificate:     certPEM,
		Signed:          signature,
	}

	stream, err := registryClient.SubscribeToStream(context.Background(), data)

	if err != nil {
		panic(err)
	}

	for {
		new, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(new)
	}
}
