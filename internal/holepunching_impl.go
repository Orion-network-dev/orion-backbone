package internal

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"golang.zx2c4.com/wireguard/wgctrl"

	"github.com/rs/zerolog/log"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	holePunchingInstances        = flag.Int("hole-punching-concurrent-instances", 16, "Concurrent hole-punching instances")
	holePunchingHost             = flag.String("hole-punching-server-address", "reg.orionet.re", "The address or name to give to the client")
	holePunchingBasePort         = flag.Int("hole-punching-base-port", 42000, "The base port for hole-punching wireguard-tunnels")
	holePunchingInterfacePrefix  = flag.String("hole-punching-interface-prefix", "reg", "The prefix added to each wireguard instance used for hole punching")
	holePunchingHandshakeTimeout = flag.Int("hole-punching-handshake-timeout-seconds", 16, "The number of seconds to wait for a client handshake")
)

type OrionHolePunchingImplementation struct {
	wgClient      *wgctrl.Client
	tasksAssigner LockableTasks
	proto.UnimplementedHolePunchingServiceServer
}

func NewOrionHolePunchingImplementation() (*OrionHolePunchingImplementation, error) {
	wg, err := wgctrl.New()
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize the wireguard control system")
		return nil, err
	}

	log.Info().Msg("initialized the Orion hole-punching api implementation")
	return &OrionHolePunchingImplementation{
		wgClient:      wg,
		tasksAssigner: NewLockableTasks(*holePunchingInstances),
	}, nil
}

func (r *OrionHolePunchingImplementation) Session(sessionInit *proto.HolePunchingInitialize, sessionServer proto.HolePunchingService_SessionServer) error {
	log.Debug().Msg("handling a hole-punching request")

	task, err := r.tasksAssigner.AssignSessionId(sessionServer.Context())
	if err != nil {
		return err
	}
	defer task.Release()

	// Parameters for the new wireguard tunnel instance used for hole-punching.
	device := wgtypes.Config{}
	port := *holePunchingBasePort + task.Id
	device.ListenPort = &port

	// Generate a preshared-key for the wireguard peer
	presharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		return err
	}

	// Add a new peer for this client
	device.Peers = append(device.Peers, wgtypes.PeerConfig{
		PublicKey:    wgtypes.Key(sessionInit.PublicKey),
		PresharedKey: &presharedKey,
	})

	// Generate a private-key for this wireguard instance
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return err
	}
	device.PrivateKey = &key

	// Create a new wireguard interface
	interfaceName := fmt.Sprintf("%s%d", *holePunchingInterfacePrefix, task.Id)
	log.Info().Str("interface-name", interfaceName).Msg("creating interface")

	wgInt, err := NewWireguardInterface(
		r.wgClient,
		&netlink.LinkAttrs{
			Name: interfaceName,
		},
		device,
	)
	if err != nil {
		return err
	}
	defer wgInt.Dispose()

	log.Debug().Msg("sending the connection information to the client")

	public_key_bytes := [wgtypes.KeyLen]byte(device.PrivateKey.PublicKey())
	preshared_key_bytes := [wgtypes.KeyLen]byte(presharedKey)

	// Send the login information the the client
	sessionServer.Send(&proto.HolePunchingEvent{
		Event: &proto.HolePunchingEvent_InitializationResponse{
			InitializationResponse: &proto.HolePunchingInitializationResponse{
				EndpointAddr: *holePunchingHost,
				EndpointPort: uint32(port),
				PublicKey:    public_key_bytes[:],
				PresharedKey: preshared_key_bytes[:],
			},
		},
	})

	// Create a new context for waiting for the first handshake from the client
	timeoutTime := time.Second * time.Duration(*holePunchingHandshakeTimeout)
	waitingCtx, ctxCancel := context.WithTimeout(sessionServer.Context(), timeoutTime)
	defer ctxCancel()
	// We're checking the status every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// We verify if an handshake was made
			log.Debug().Str("interface", interfaceName).Msg("checking the wireguard interface for handshakes")
			dev, err := r.wgClient.Device(interfaceName)
			if err != nil {
				log.Error().Err(err).Msg("error while reading the interface information")
				return err
			}

			if len(dev.Peers) != 1 {
				err = fmt.Errorf("more than one peer is connected to the hole-punching instance")
				log.Error().Err(err).Msg("this should be not possible")
				return err
			}

			peer := dev.Peers[0]
			// We check if an endpoint was recorded
			if peer.Endpoint != nil {
				log.Debug().Int("task-id", task.Id).IPAddr("address", peer.Endpoint.IP).Int("port", peer.Endpoint.Port).Msg("got a connection to the wireguard instance")

				sessionServer.Send(&proto.HolePunchingEvent{
					Event: &proto.HolePunchingEvent_Complete{
						Complete: &proto.HolePunchingCompleteResponse{
							ClientEndpointAddr: peer.Endpoint.IP.String(),
							ClientEndpointPort: uint32(peer.Endpoint.Port),
						},
					},
				})

				return nil
			}

		case <-waitingCtx.Done():
			return waitingCtx.Err()
		}
	}
}
