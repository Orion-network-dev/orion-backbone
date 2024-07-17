package implementation

import (
	"context"
	"fmt"
	"net"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// This is called wwhen another client wants to connect to the current client
// This causes the current node to do a few things:
//   - Initialize a new wireguard tunnel
//   - Do hole-punching to expose a public port and address to the peer
//   - Set a timeout of 5minutes for the client to connect to the endpoint
//     a) If the client failed to connect, we simply discard the wireguard tunnel
//     b) Is the client succesfullt connected, we do a test ping and setup monitoring
func (c *OrionClientDaemon) handleWantsToConnect(event *proto.ClientWantToConnectToClient) {
	// Check if this is a valid request
	if event.DestinationPeerId != c.memberId || event.SourcePeerId == c.memberId {
		return
	}

	// We generate a new key-pair for the new tunnel
	privatekey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate the private / public key pair")
		return
	}
	publickey := privatekey.PublicKey()

	// All orion connections use a preshared key to ensure Quantic-safe tunnels
	// This only appends on the responding side.
	presharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate the preshared key")
		return
	}

	// We create the new wireguard tunnel
	tunnel, err := internal.NewWireguardInterface(c.wgClient, &netlink.LinkAttrs{
		Name:  fmt.Sprintf("orion%d", event.SourcePeerId),
		Group: 30,
	}, wgtypes.Config{
		PrivateKey: &privatekey,
	})
	if err != nil {
		log.Error().Err(err).Msg("cannot make wireguard interface")
		return
	}
	c.wireguardTunnels[event.SourcePeerId] = tunnel

	// We need to do holepunching first because hole-punching overrides the tunnel's peers
	holePunchingResult, err := holePunchTunnel(context.Background(), c.wgClient, tunnel, c.holePunchingClient)
	if err != nil {
		log.Error().Err(err).Msg("cannot holepunch interface")
		return
	}

	// Calculate the ip address
	selfIP, otherIP, err := internal.GetSelfAddress(uint32(*memberId), event.SourcePeerId)
	if err != nil {
		log.Error().Err(err).Msg("failed to calculate the IP adresses from the adressing plan")
		return
	}

	tunnel.SetPeers(c.wgClient, []wgtypes.PeerConfig{
		{
			Endpoint: &net.UDPAddr{
				IP:   net.ParseIP(event.EndpointAddr),
				Port: int(event.EndpointPort),
			},
			PresharedKey:                &presharedKey,
			PublicKey:                   wgtypes.Key(event.PublicKey),
			PersistentKeepaliveInterval: keepAlive,
			AllowedIPs: []net.IPNet{
				allIPRanges,
			},
		},
	})
	tunnel.SetAddress(selfIP)

	err = c.registryStream.Send(&proto.RPCClientEvent{
		Event: &proto.RPCClientEvent_ConnectResponse{
			ConnectResponse: &proto.ClientWantToConnectToClientResponse{
				EndpointAddr:      holePunchingResult.ClientEndpointAddr,
				EndpointPort:      holePunchingResult.ClientEndpointPort,
				PublicKey:         publickey[:],
				FriendlyName:      *friendlyName,
				SourcePeerId:      c.memberId,
				DestinationPeerId: event.SourcePeerId,
				PresharedKey:      presharedKey[:],
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to send the message to the server")
		return
	}

	c.frrManager.Peers[event.SourcePeerId] = &Peer{
		ASN:     event.SourcePeerId + 64511,
		Address: otherIP.IP.String(),
	}
	err = c.frrManager.Update()

	if err != nil {
		log.Error().Err(err).Msg("failed to apply the modifications to FRR's BGPD")
	}
}
