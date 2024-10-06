package link

import (
	"fmt"
	"math"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal"
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

		// todo: termination of connection

		latency = time.Hour * 24 * 7
	}

	log.Debug().Dur("ping-reponse", latency).Msg("ping(ed) peer")
	newPeer := c.frrManager.GetPeer(internal.IdentityFromRouter(c.other))
	if newPeer == nil {
		return fmt.Errorf("peer is not existant")
	}

	// f\left(x\right)=\min\left(\max\left(e^{\ \left(\frac{500-x}{80}\right)},0\right),300\right)
	newPeer.Weight = uint32(math.Min(
		300,
		math.Max(
			math.Exp(
				(500-float64(latency.Milliseconds()))/80,
			),
			0,
		),
	))
	c.frrManager.UpdatePeer(internal.IdentityFromRouter(c.other), newPeer)

	err = c.frrManager.Update()
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
					Uint64("peer-id", internal.IdentityFromRouter(c.other)).
					Msgf("failed to adjust the weight")
			}

		case <-c.ctx.Done():
			log.Error().
				Err(c.ctx.Err()).
				Uint64("peer-id", internal.IdentityFromRouter(c.other)).
				Msgf("ending the background task")
			return
		}
	}
}
