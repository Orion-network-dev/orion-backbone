package implementation

import (
	"context"
	"fmt"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (c *OrionClientDaemon) handleNewClient(event *proto.ClientNewOnNetworkEvent) {
	log.Debug().Msg("got new client message, trying to initialize a p2p connection")

	privatekey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		log.Error().Err(err).Msg("failure to generate a wireguard private key")
		return
	}
	publickey := privatekey.PublicKey()

	tunnel, err := internal.NewWireguardInterface(c.wgClient, &netlink.LinkAttrs{
		Name:  fmt.Sprintf("orion%d", event.PeerId),
		Group: 30,
	}, wgtypes.Config{
		PrivateKey:   &privatekey,
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{},
	})
	if err != nil {
		log.Error().Err(err).Msg("cannot make wireguard interface")
		return
	}
	c.wireguardTunnels[event.PeerId] = tunnel

	ctx := context.Background()
	holepunch, err := holePunchTunnel(ctx, c.wgClient, tunnel, c.holePunchingClient)
	if err != nil {
		log.Error().Err(err).Msg("cannot hole punch interface")
		tunnel.Dispose()
		return
	}

	// Ask a new connection by emitting a client event
	err = c.registryStream.Send(&proto.RPCClientEvent{
		Event: &proto.RPCClientEvent_Connect{
			Connect: &proto.ClientWantToConnectToClient{
				EndpointAddr:      holepunch.ClientEndpointAddr,
				EndpointPort:      holepunch.ClientEndpointPort,
				PublicKey:         publickey[:],
				FriendlyName:      *friendlyName,
				DestinationPeerId: event.PeerId,
				SourcePeerId:      c.memberId,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msgf("couldn't swrite the initialization message to the gRPC connection")
	}

}
