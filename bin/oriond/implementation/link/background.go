package link

import (
	"math"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation/frr"
	"github.com/go-ping/ping"
)

func (c *PeerLink) updateWeights() error {
	pinger, err := ping.NewPinger(c.otherIP.IP.String())
	if err != nil {
		return err
	}
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
	weight := uint32(math.Min(math.Exp(float64(latency.Milliseconds()/15))+50, 0))
	c.frrManager.Peers[c.otherID] = &frr.Peer{
		ASN:     c.otherID + 64511,
		Address: c.otherIP.IP.String(),
		Weight:  weight,
	}
	c.frrManager.Update()

	return nil
}

func (c *PeerLink) backgroundTask() {
	// We check the status every 60 seconds
	timer := time.NewTicker(time.Second * 60)

	for {
		select {
		case <-timer.C:
			c.updateWeights()

		case <-c.ctx.Done():
			return
		}
	}
}
