package link

import (
	"math"
	"time"

	"github.com/go-ping/ping"
	"github.com/rs/zerolog/log"
)

func (c *PeerLink) updateWeights() error {
	log.Debug().Msg("starting to ping")
	pinger, err := ping.NewPinger(c.otherIP.IP.String())
	if err != nil {
		return err
	}
	pinger.SetPrivileged(true)

	// Ping one time
	pinger.Count = 1
	pinger.Timeout = time.Second * 5

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

	// f\left(x\right)=\min\left(\max\left(e^{\ \left(\frac{500-x}{80}\right)},0\right),300\right)
	metric := int(math.Min(
		300,
		math.Max(
			math.Exp(
				(500-float64(latency.Milliseconds()))/80,
			),
			0,
		),
	))
	err = c.wireguardTunnel.SetMetric(metric)
	if err != nil {
		return err
	}

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
