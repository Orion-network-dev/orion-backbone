package implementation

import (
	"net"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog/log"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (c *OrionClientDaemon) handleWantsToConnectResponse(event *proto.ClientWantToConnectToClientResponse) {
	wg := c.wireguardTunnels[event.SourcePeerId]

	if wg != nil {
		// Calculate the ip address
		selfIP, otherIP, err := internal.GetSelfAddress(c.memberId, event.SourcePeerId)
		if err != nil {
			log.Error().Err(err).Msg("failed to compute the self address")
			return
		}

		wg.SetPeers(c.wgClient, []wgtypes.PeerConfig{
			{
				Endpoint: &net.UDPAddr{
					IP:   net.ParseIP(event.EndpointAddr),
					Port: int(event.EndpointPort),
				},
				PresharedKey:                (*wgtypes.Key)(event.PresharedKey),
				PublicKey:                   wgtypes.Key(event.PublicKey),
				PersistentKeepaliveInterval: keepAlive,
				AllowedIPs: []net.IPNet{
					allIPRanges,
				},
			},
		})
		wg.SetAddress(selfIP)

		c.frrManager.Peers[event.SourcePeerId] = &Peer{
			ASN:     event.SourcePeerId + 64511,
			Address: otherIP.IP.String(),
		}
		err = c.frrManager.Update()
		if err != nil {
			log.Error().Err(err).Msg("failed to update the frr configuration")
			return
		}
	}

}
