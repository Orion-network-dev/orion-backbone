package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"net"
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
	debug          = flag.Bool("debug", false, "change the log level to debug")
	friendlyName   = flag.String("friendly-name", "", "the public friendly name the instance will have")
	memberId       = flag.Int("member-id", 0, "the public friendly name the instance will have")
	registryServer = flag.String("registry-server", "puffer.fish", "the address of the registry server")
	registryPort   = flag.Uint("registry-port", 6443, "the port used by the registry")
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

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", *registryServer, *registryPort), grpc.WithTransportCredentials(cred), grpc.WithIdleTimeout(time.Second*120))
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

	tunnels := make([]*internal.WireguardInterface, 255)
	defer func() {
		for _, tunnel := range tunnels {
			tunnel.Dispose()
		}
	}()

	// Go routine used to login
	go func() {
		log.Debug().Msg("preparing to send the initialization message for authentication")

		// Reading
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

		err = stream.Send(&proto.RPCClientEvent{
			Event: &proto.RPCClientEvent_Initialize{
				Initialize: internal.CalculateNonce(int64(*memberId), *friendlyName, certPEM, pk),
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
		minute := time.Minute

		if new_client := data.GetNewClient(); new_client != nil {
			log.Debug().Msg("got new client message, trying to initialize a p2p connection")

			privatekey, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				log.Fatal().Err(err).Msg("failure to generate a wireguard private key")
			}
			publickey := privatekey.PublicKey()

			tunnel, err := internal.NewWireguardInterface(wgClient, &netlink.LinkAttrs{
				Name: fmt.Sprintf("orion%d", new_client.PeerId),
			}, wgtypes.Config{
				PrivateKey:   &privatekey,
				ReplacePeers: true,
				Peers:        []wgtypes.PeerConfig{},
			})
			if err != nil {
				log.Fatal().Err(err).Msg("cannot make wireguard interface")
			}
			tunnels[new_client.PeerId] = tunnel

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
				log.Fatal().Err(err).Msg("cannot make wireguard interface")
			}
			publickey := privatekey.PublicKey()
			presharedKey, err := wgtypes.GenerateKey()
			if err != nil {
				log.Fatal().Err(err).Msg("cannot make wireguard interface")
			}
			tunnel, err := internal.NewWireguardInterface(wgClient, &netlink.LinkAttrs{
				Name: fmt.Sprintf("orion%d", wants_to.SourcePeerId),
			}, wgtypes.Config{
				PrivateKey: &privatekey,
			})
			if err != nil {
				log.Fatal().Err(err).Msg("cannot make wireguard interface")
			}
			tunnels[wants_to.SourcePeerId] = tunnel
			result, err := internal.HolePunchTunnel(context.Background(), wgClient, tunnel, holepunchingClient)
			if err != nil {
				log.Fatal().Err(err).Msg("cannot holepunch interface")
			}
			tunnel.SetPeers(wgClient, []wgtypes.PeerConfig{
				{
					Endpoint: &net.UDPAddr{
						IP:   net.ParseIP(wants_to.EndpointAddr),
						Port: int(wants_to.EndpointPort),
					},
					PresharedKey:                &presharedKey,
					PublicKey:                   wgtypes.Key(wants_to.PublicKey),
					PersistentKeepaliveInterval: &minute,
				},
			})

			response := &proto.ClientWantToConnectToClientResponse{
				EndpointAddr:      result.ClientEndpointAddr,
				EndpointPort:      result.ClientEndpointPort,
				PublicKey:         publickey[:],
				FriendlyName:      *friendlyName,
				SourcePeerId:      int64(*memberId),
				DestinationPeerId: wants_to.SourcePeerId,
				PresharedKey:      presharedKey[:],
			}
			fmt.Println(response)
			stream.Send(&proto.RPCClientEvent{
				Event: &proto.RPCClientEvent_ConnectResponse{
					ConnectResponse: response,
				},
			})
			continue
		}

		if response := data.GetWantsToConnectResponse(); response != nil {
			// Now that the connection is done, we simply need to add the peer
			wg := tunnels[response.SourcePeerId]
			if wg != nil {
				wg.SetPeers(wgClient, []wgtypes.PeerConfig{
					{
						Endpoint: &net.UDPAddr{
							IP:   net.ParseIP(response.EndpointAddr),
							Port: int(response.EndpointPort),
						},
						PresharedKey:                (*wgtypes.Key)(response.PresharedKey),
						PublicKey:                   wgtypes.Key(response.PublicKey),
						PersistentKeepaliveInterval: &minute,
					},
				})
			}
		}
	}
}
