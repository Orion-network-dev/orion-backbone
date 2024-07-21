package link

func (c *PeerLink) Dispose() error {
	// we remove the frr peer
	c.frrManager.Peers[c.otherID] = nil
	err := c.frrManager.Update()
	if err != nil {
		return err
	}

	// we dispose the vpn tunnel
	c.wireguardTunnel.Dispose()
	// cancel all the running tasks with the peer's context
	c.cancel()

	return nil
}
