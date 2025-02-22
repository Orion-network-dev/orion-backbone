package link

import (
	"math"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/rs/zerolog/log"
)

func (c *PeerLink) updateWeights() error {
	log.Debug().Msg("starting to ping")
	pinger, err := probing.NewPinger(c.otherIP.IP.String())
	if err != nil {
		return err
	}
	pinger.SetPrivileged(true)

	// Ping one time
	pinger.Count = 1
	pinger.Timeout = time.Second * 5
	pinger.InterfaceName = c.wireguardTunnel.WgLink.InterfaceAttrs.Name

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return err
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats

	latency := stats.AvgRtt
	if stats.PacketLoss > 0 {
		log.Debug().Msg("ping failed")
		latency = time.Hour * 24 * 7
	}

	log.Debug().Dur("ping-reponse", latency).Msg("ping(ed) peer")

	// max( 0, latency ), more when less latency & less with more latency
	metric := math.Max(float64(0), float64(int64(500)-latency.Milliseconds()))
	peer := c.frrManager.GetPeer(c.otherID)
	peer.Weight = uint32(metric)
	// c.frrManager.UpdatePeer(c.otherID, peer)
	// c.frrManager.Update()

	return nil
}

func (c *PeerLink) backgroundTask() {
	// We check the status every 60 seconds
	timer := time.NewTicker(time.Second * 60)

	for {
		select {
		case <-timer.C:
			if err := c.updateWeights(); err != nil {
				log.Error().
					Err(err).
					Uint32("peer-id", c.otherID).
					Msgf("failed to adjust the weight")
			}

		case <-c.ctx.Done():
			log.Error().
				Err(c.ctx.Err()).
				Uint32("peer-id", c.otherID).
				Msgf("ending the background task")
			return
		}
	}
}
