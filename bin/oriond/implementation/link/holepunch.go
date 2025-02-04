package link

import (
	"context"
	"flag"
	"fmt"

	"github.com/orion-network-dev/orion-backbone/internal"
	"github.com/orion-network-dev/orion-backbone/internal/state"
	"github.com/rs/zerolog/log"
	"gitlab.com/NebulousLabs/go-upnp"
)

var (
	holePunchOverrideAddress = flag.String("override-hole-punch-address", "", "Override the public port for this instance")
	holePunchUPNPEnable      = flag.Bool("hole-punch-upnp", true, "Specify is uPNP should be enabled")
)

func (c *PeerLink) upnpInit() {
	if *holePunchUPNPEnable {
		digd, err := upnp.Discover()
		if err != nil {
			log.Err(err).Msg("failed to intitialize upnp")
			return
		}

		ip, err := digd.ExternalIP()
		if err != nil {
			log.Err(err).Msg("failed to intitialize upnp")
			return
		}

		log.Info().Msgf("upnp status initialized")
		c.externalIP = &ip
		c.igd = digd
	}
}

func (c *PeerLink) HolePunchTunnel(
	parentCtx context.Context,
	tunnel *internal.WireguardInterface,
) *state.Endpoint {
	if c.igd != nil {
		log.Debug().Msg("trying to open port using upnp")
		forwarded, err := c.igd.IsForwardedUDP(uint16(*tunnel.WgConfig.ListenPort))
		if err != nil {
			log.Err(err).Msg("failed to open port using upnp")
		} else {
			if forwarded {
				log.Debug().Msg("port already openned using upnp")
				return &state.Endpoint{
					Address:    *c.externalIP,
					PublicPort: uint16(*tunnel.WgConfig.ListenPort),
					PublicKey:  tunnel.WgConfig.PrivateKey.PublicKey().String(),
				}
			}

			if err := c.igd.ForwardUDP(uint16(*tunnel.WgConfig.ListenPort), fmt.Sprintf("Used by orion to peer %d", c.otherID)); err != nil {
				log.Debug().Msg("failed to pen port using upnp")
			}
		}
	}

	return nil
}
