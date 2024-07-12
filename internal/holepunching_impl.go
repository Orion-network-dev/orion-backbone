package internal

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"golang.zx2c4.com/wireguard/wgctrl"

	"github.com/rs/zerolog/log"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type OrionHolePunchingImplementation struct {
	wgClient      *wgctrl.Client
	tasksAssigner LockableTasks
	proto.UnimplementedHolePunchingServiceServer
}

func NewOrionHolePunchingImplementations() (*OrionHolePunchingImplementation, error) {
	wg, err := wgctrl.New()
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize the wireguard control system")
		return nil, err
	}

	log.Info().Msg("initialized the Orion hole-punching api implementation")
	return &OrionHolePunchingImplementation{
		wgClient:      wg,
		tasksAssigner: NewLockableTasks(255),
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

	// Generate a new preshared key for this link
	presharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		return err
	}

	// Add a new peer for the client.
	device.Peers = append(device.Peers, wgtypes.PeerConfig{
		PublicKey:    wgtypes.Key(sessionInit.PublicKey),
		PresharedKey: &presharedKey,
		AllowedIPs: []net.IPNet{
			{
				IP:   net.IPv4(10, 255, byte(task.Id), 0),
				Mask: net.CIDRMask(31, 32),
			},
		},
	})

	// Specifying a new port
	port := 42000 + task.Id
	device.ListenPort = &port

	// Generating a new private key for our tunnel.
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return err
	}
	device.PrivateKey = &key
	int_name := fmt.Sprintf("reg%d", task.Id)
	log.Info().Str("interface-name", int_name).Msg("creating interface")

	// Creating link using the `netlink` package
	wglink := WireguardNetLink{
		Id:     task.Id,
		Prefix: "reg",
	}
	err = netlink.LinkAdd(wglink)
	if err != nil {
		log.Error().Err(err).Msg("error while creating the interface")
		return err
	}
	defer netlink.LinkDel(wglink)

	// Configuring the device using our instance
	err = r.wgClient.ConfigureDevice(int_name, device)
	if err != nil {
		log.Error().Err(err).Msg("failed to apply the wireguard configuration")
		return err
	}

	ipConfig := &netlink.Addr{IPNet: &net.IPNet{
		IP:   net.IPv4(10, 255, byte(task.Id), 1),
		Mask: net.CIDRMask(24, 32),
	}}

	if err = netlink.AddrAdd(wglink, ipConfig); err != nil {
		log.Error().Err(err).Msg("failed to add the ip configuration")
		return err
	}
	if err = netlink.LinkSetUp(wglink); err != nil {
		log.Error().Err(err).Msg("failed to set the interface up")
		return err
	}

	log.Debug().Msg("sending the connection information to the client")

	publick := [wgtypes.KeyLen]byte(device.PrivateKey.PublicKey())
	presharedk := [wgtypes.KeyLen]byte(presharedKey)
	// Sending the connection informations to the client.
	sessionServer.Send(&proto.HolePunchingEvent{
		Event: &proto.HolePunchingEvent_InitializationResponse{
			InitializationResponse: &proto.HolePunchingInitializationResponse{
				EndpointAddr:  "reg.orionet.re",
				EndpointPort:  uint32(port),
				PublicKey:     publick[:],
				PresharedKey:  presharedk[:],
				ClientAddress: fmt.Sprintf("10.255.%d.2", task.Id),
				RemoteAddress: fmt.Sprintf("10.255.%d.1", task.Id),
			},
		},
	})

	waitingCtx, ctxCancel := context.WithTimeout(sessionServer.Context(), time.Second*15)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Debug().Str("interface", int_name).Msg("ticking the interface")
			dev, err := r.wgClient.Device(int_name)

			if err != nil {
				log.Error().Err(err).Msg("error while reading the interface information")
				ctxCancel()
				break
			}
			if len(dev.Peers) != 1 {
				log.Error().Msg("more than one peer is connecte to the hole-punching instance")
				ctxCancel()
				break
			}

			peer := dev.Peers[0]
			if peer.Endpoint != nil {
				log.Info().Int("task-id", task.Id).IPAddr("address", peer.Endpoint.IP).Int("port", peer.Endpoint.Port).Msg("got a connection to the wireguard instance")
				sessionServer.Send(&proto.HolePunchingEvent{
					Event: &proto.HolePunchingEvent_Complete{
						Complete: &proto.HolePunchingCompleteResponse{
							ClientEndpoint: fmt.Sprintf("%s:%d", peer.Endpoint.IP, peer.Endpoint.Port),
						},
					},
				})
				ctxCancel()
				return nil
			}

		case <-waitingCtx.Done():
			ctxCancel()
			return waitingCtx.Err()
		}
	}
}
