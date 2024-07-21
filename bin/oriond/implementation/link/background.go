package link

import (
	"math"
	"time"

	"github.com/go-ping/ping"
	"github.com/rs/zerolog/log"
)

func (c *PeerLink) updateWeights() error {
	pinger, err := ping.NewPinger(c.otherIP.IP.String())
	if err != nil {
		return err
	}
	pinger.SetPrivileged(true)

	// Ping one time
	pinger.Count = 1
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return err
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	latency := stats.AvgRtt
	if stats.PacketsRecv == stats.PacketsSent {
		latency = time.Hour * 24 * 7
		return nil
	}

	// f\left(x\right)=-e^{\ \left(\frac{x}{15}\right)}+50
	c.frrManager.Peers[c.otherID].Weight = uint32(math.Min(math.Exp(float64(latency.Milliseconds()/15))+50, 0))
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
