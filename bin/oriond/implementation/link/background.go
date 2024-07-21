package link

import "time"

func (c *PeerLink) updateWeights() error {
	// todo: do an icmp ping and ajust the bgp weights related to it
	// latency := 0

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
