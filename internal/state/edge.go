package state

type Edge struct {
	RouterA *Router
	RouterB *Router
}

// Sends a new initialization step to both the peers
// asking for a random one to choose a pre-shared key
// The hole-punching logic is done locally by the peers
func (c *Edge) Initialize() {

	go func() {
		select {
		case <-c.RouterA.routerObjectContext.Done():
		case <-c.RouterB.routerObjectContext.Done():
		}

		// the edge lifetime is finished
		// todo: teardown edge
	}()
}
