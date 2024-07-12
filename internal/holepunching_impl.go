package internal

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"golang.zx2c4.com/wireguard/wgctrl"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type OrionHolePunchingImplementations struct {
	WgClient *wgctrl.Client
	proto.UnimplementedHolePunchingServiceServer
}

type WireguardNetLink struct {
	netlink.Link
	Id     int
	Prefix string
}

func (r WireguardNetLink) Type() string {
	return "wireguard"
}

func (r WireguardNetLink) Attrs() *netlink.LinkAttrs {
	return &netlink.LinkAttrs{
		Name: fmt.Sprintf("%s%d", r.Prefix, r.Id),
	}
}

func (r OrionHolePunchingImplementations) Session(sessionInit *proto.HolePunchingInitialize, sessionServer proto.HolePunchingService_SessionServer) error {
	// Generate a new preshared key for this link
	presharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		return err
	}

	// Generating an id for our client.
	id := rand.Int() % 255

	// Parameters for the new wireguard tunnel instance used for hole-punching.
	device := wgtypes.Config{}

	// Add a new peer for the client.
	device.Peers = append(device.Peers, wgtypes.PeerConfig{
		PublicKey:    wgtypes.Key(sessionInit.PublicKey),
		PresharedKey: &presharedKey,
		AllowedIPs: []net.IPNet{
			{
				IP:   net.IPv4(10, 255, byte(id), 0),
				Mask: net.CIDRMask(31, 32),
			},
		},
	})

	// Specify that we want to replace all the existing peers.
	device.ReplacePeers = false
	// Specifying a new port
	port := 42000 + id
	device.ListenPort = &port

	// Generating a new private key for our tunnel.
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return err
	}
	device.PrivateKey = &key
	int_name := fmt.Sprintf("reg%d", id)
	fmt.Printf("Creating %s\n", int_name)

	// Creating link using the `netlink` package
	wglink := WireguardNetLink{
		Id:     id,
		Prefix: "reg",
	}
	err = netlink.LinkAdd(wglink)
	if err != nil {
		return err
	}

	defer func() {
		netlink.LinkDel(wglink)
	}()

	// Configuring the device using our instance
	err = r.WgClient.ConfigureDevice(int_name, device)
	if err != nil {
		return err
	}

	// Set the server IP address on the tunnel
	lnk, err := netlink.LinkByName(int_name)
	if err != nil {
		return err
	}

	ipConfig := &netlink.Addr{IPNet: &net.IPNet{
		IP:   net.IPv4(10, 255, byte(id), 1),
		Mask: net.CIDRMask(24, 32),
	}}

	if err = netlink.AddrAdd(lnk, ipConfig); err != nil {
		return err
	}
	if err = netlink.LinkSetUp(lnk); err != nil {
		return err
	}

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
				ClientAddress: fmt.Sprintf("10.255.%d.2", id),
				RemoteAddress: fmt.Sprintf("10.255.%d.1", id),
			},
		},
	})

	waitingCtx, ctxCancel := context.WithTimeout(sessionServer.Context(), time.Second*30)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Get the peer
			dev, err := r.WgClient.Device(int_name)

			if err != nil {
				ctxCancel()
				break
			}
			if len(dev.Peers) != 1 {
				ctxCancel()
				break
			}

			peer := dev.Peers[0]
			if peer.Endpoint != nil {
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
