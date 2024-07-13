package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/grpc"
)

var (
	debug        = flag.Bool("debug", false, "change the log level to debug")
	friendlyName = flag.String("friendly-name", "", "the public friendly name the instance will have")
	memberId     = flag.Int("member-id", 0, "the public friendly name the instance will have")
)

func main() {
	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Get TLS credentials
	cred, err := internal.LoadTLS(false)
	if err != nil {
		log.Fatal().Msgf("Unable to connect gRPC channel %v", err)
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", "reg.orionet.re", 6443), grpc.WithTransportCredentials(cred), grpc.WithIdleTimeout(time.Second*120))
	if err != nil {
		log.Fatal().Msgf("Unable to connect gRPC channel %v", err)
	}

	// Create the gRPC client
	registryClient := proto.NewRegistryClient(conn)
	// holepunchingClient = proto.NewHolePunchingServiceClient(conn)

	stream, err := registryClient.SubscribeToStream(context.Background())
	if err != nil {
		panic(err)
	}

	// Go routine used to login
	go func() {
		time := time.Now().Unix()
		nonce := sha512.New().Sum([]byte(internal.CalculateNonce(int64(*memberId), *friendlyName, time)))
		certPEM, err := os.ReadFile(*internal.CertificatePath)
		if err != nil {
			panic(err)
		}
		privateKey, err := os.ReadFile(*internal.KeyPath)
		if err != nil {
			panic(err)
		}
		zzz, _ := pem.Decode(privateKey)

		pk, err := x509.ParseECPrivateKey(zzz.Bytes)
		if err != nil {
			panic(err)
		}
		signature, err := ecdsa.SignASN1(rand.Reader, pk, nonce)
		if err != nil {
			panic(err)
		}
		stream.Send(&proto.RPCClientEvent{
			Event: &proto.RPCClientEvent_Initialize{
				Initialize: &proto.InitializeRequest{
					FriendlyName:    *friendlyName,
					TimestampSigned: time,
					MemberId:        int64(*memberId),
					Certificate:     certPEM,
					Signed:          signature,
				},
			},
		})
	}()

	for {
		data, err := stream.Recv()
		if err != nil {
			break
		}

		if new_client := data.GetNewClient(); new_client != nil {
			privatekey, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				panic(err)
			}
			publickey := privatekey.PublicKey()
			presharedkey, err := wgtypes.GenerateKey()
			if err != nil {
				panic(err)
			}

			// Ask a new connection by emitting a client event
			stream.Send(&proto.RPCClientEvent{
				Event: &proto.RPCClientEvent_Connect{
					Connect: &proto.ClientWantToConnectToClient{
						EndpointAddr:      "127.0.0.1",
						EndpointPort:      5001,
						PublicKey:         publickey[:],
						PresharedKey:      presharedkey[:],
						FriendlyName:      *friendlyName,
						DestinationPeerId: new_client.PeerId,
						SourcePeerId:      int64(*memberId),
					},
				},
			})
		}

		if wants_to := data.GetWantsToConnect(); wants_to != nil {
			// todo: create interface
			privatekey, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				panic(err)
			}
			publickey := privatekey.PublicKey()

			response := &proto.ClientWantToConnectToClientResponse{
				EndpointAddr:      "127.0.0.1",
				EndpointPort:      5001,
				PublicKey:         publickey[:],
				FriendlyName:      *friendlyName,
				SourcePeerId:      int64(*memberId),
				DestinationPeerId: wants_to.SourcePeerId,
			}
			fmt.Println(response)
			stream.Send(&proto.RPCClientEvent{
				Event: &proto.RPCClientEvent_ConnectResponse{
					ConnectResponse: response,
				},
			})

			fmt.Println("Client connection in progress.")
		}

		if response := data.GetWantsToConnectResponse(); response != nil {
			fmt.Println("Peer responded.")
			fmt.Print(response)
		}
	}
}
