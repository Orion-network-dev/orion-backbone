package link

import "github.com/MatthieuCoder/OrionV3/internal"

func (c *PeerLink) Dispose() error {
	// cancel all the running tasks with the peer's context
	c.cancel()
	// we remove the frr peer
	c.frrManager.UpdatePeer(internal.IdentityFromRouter(c.other), nil)
	err := c.frrManager.Update()
	if err != nil {
		return err
	}

	// we dispose the vpn tunnel
	c.wireguardTunnel.Dispose()

	return nil
}

func (c *PeerLink) Terminate() {

}
