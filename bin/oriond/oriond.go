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
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
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

	wgClient, err := wgctrl.New()
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to connect to wireguard")
	}

	// Get TLS credentials
	cred, err := internal.LoadTLS(false)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to connect gRPC channel")
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", "reg.orionet.re", 6443), grpc.WithTransportCredentials(cred), grpc.WithIdleTimeout(time.Second*120))
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to connect gRPC channel")
	}

	// Create the gRPC client
	registryClient := proto.NewRegistryClient(conn)
	holepunchingClient := proto.NewHolePunchingServiceClient(conn)

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
			log.Fatal().Err(err).Msgf("coundn't open the certificate pem file")
		}

		privateKey, err := os.ReadFile(*internal.KeyPath)
		if err != nil {
			log.Fatal().Err(err).Msgf("coundn't open the certificate key file")
		}
		rawCertificate, _ := pem.Decode(privateKey)
		pk, err := x509.ParseECPrivateKey(rawCertificate.Bytes)

		if err != nil {
			log.Fatal().Err(err).Msgf("coundn't read the certificate key file")
		}

		signature, err := ecdsa.SignASN1(rand.Reader, pk, nonce)
		if err != nil {
			log.Fatal().Err(err).Msgf("couldn't sign the nonce data")
		}

		err = stream.Send(&proto.RPCClientEvent{
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

		if err != nil {
			log.Fatal().Err(err).Msgf("couldn't swrite the initialization message to the gRPC connection")
		}
	}()

	for {
		data, err := stream.Recv()
		if err != nil {
			log.Fatal().Err(err).Msg("failure while reading the gRPC stream")
		}

		if new_client := data.GetNewClient(); new_client != nil {
			log.Debug().Msg("got new client message, trying to initialize a p2p connection")

			privatekey, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				log.Fatal().Err(err).Msg("failure to generate a wireguard private key")
			}
			publickey := privatekey.PublicKey()

			tunnel, err := internal.NewWireguardInterface(wgClient, &netlink.LinkAttrs{
				Name: fmt.Sprintf("peer%d", new_client.PeerId),
			}, wgtypes.Config{
				PrivateKey:   &privatekey,
				ReplacePeers: true,
				Peers:        []wgtypes.PeerConfig{},
			})
			if err != nil {
				log.Fatal().Err(err).Msg("cannot make wireguard interface")
			}

			ctx := context.Background()
			holepunch, err := internal.HolePunchTunnel(ctx, wgClient, tunnel, holepunchingClient)
			if err != nil {
				log.Error().Err(err).Msg("cannot hole punch interface")
				tunnel.Dispose()
				continue
			}

			// Ask a new connection by emitting a client event
			err = stream.Send(&proto.RPCClientEvent{
				Event: &proto.RPCClientEvent_Connect{
					Connect: &proto.ClientWantToConnectToClient{
						EndpointAddr:      holepunch.ClientEndpointAddr,
						EndpointPort:      holepunch.ClientEndpointPort,
						PublicKey:         publickey[:],
						FriendlyName:      *friendlyName,
						DestinationPeerId: new_client.PeerId,
						SourcePeerId:      int64(*memberId),
					},
				},
			})
			if err != nil {
				log.Fatal().Err(err).Msgf("couldn't swrite the initialization message to the gRPC connection")
			}

			continue
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
			continue
		}
		fmt.Println(data)

		if response := data.GetWantsToConnectResponse(); response != nil {
			fmt.Println("Peer responded.")
		}
	}
}
