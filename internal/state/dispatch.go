package state

import (
	"slices"
)

func (c *OrionRegistryState) GetRouters() OrionRegistryRoutersState {
	return c.routers
}

// Called once a new member is joining the network
func (c *OrionRegistryState) DispatchNewRouterEvent(
	newRouter *Router,
) {
	c.routersLock.Lock()
	defer c.routersLock.Unlock()

	c.routers[RouterIdentity(newRouter.Identity)] = newRouter

	for _, router := range c.routers {
		if router.Identity != newRouter.Identity {
			router.DispatchNewRouterEvent(newRouter)
		}
	}
}

// Called once a member is removed
func (c *OrionRegistryState) DispatchRouterRemovedEvent(
	deletedRouter *Router,
) {
	if c.routers[deletedRouter.Identity] == nil {
		return
	}

	c.routersLock.Lock()
	defer c.routersLock.Unlock()

	for _, router := range c.routers {
		// todo: send event to other peers
		if router.Identity != deletedRouter.Identity {
			// todo: send event
			router.DispatchRouterRemovedEvent(deletedRouter)
		}
	}

	c.routers[deletedRouter.Identity].dispose()
	c.routers[deletedRouter.Identity] = nil
}

// Dispatch new connection
func (c *OrionRegistryState) DispatchNewEdge(
	edge *Edge,
) {
	c.edgesLock.Lock()
	defer c.edgesLock.Unlock()

	edgeA := edge.RouterA
	edgeB := edge.RouterB

	// check that the routers are existing
	if edgeA == nil || edgeB == nil {
		// todo: report error
		return
	}

	ids := []uint32{uint32(edgeA.Identity), uint32(edgeB.Identity)}
	slices.Sort(ids)

	// concatenate the bits
	edgeId := uint64(ids[0]) << 32 & uint64(ids[1])

	c.edges[edgeId] = edge
	//edgeA.NewEdge(edge)
	//edgeB.NewEdge(edge)

}
