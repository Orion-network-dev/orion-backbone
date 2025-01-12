package state

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type OrionRegistryRoutersState map[RouterIdentity]*Router
type OrionRegistryEdgesState map[EdgeIdentity]*Edge

// State related to the Orion-Registry
// component, that handles all the connection
// initialization system.
type OrionRegistryState struct {
	routers     OrionRegistryRoutersState // List of routers
	routersLock sync.Mutex

	edges     OrionRegistryEdgesState // List of edges in the orion graph
	edgesLock sync.Mutex
}

func NewOrionRegistryState() *OrionRegistryState {
	return &OrionRegistryState{
		routers:     make(OrionRegistryRoutersState),
		routersLock: sync.Mutex{},
		edges:       make(OrionRegistryEdgesState),
		edgesLock:   sync.Mutex{},
	}
}

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
			router.Dispatch(RouterConnectEvent{
				Router: newRouter,
			})
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

	c.routers[deletedRouter.Identity].dispose()
	c.routers[deletedRouter.Identity] = nil
}

// Dispatch new connection
func (c *OrionRegistryState) DispatchNewEdge(
	edge *Edge,
) {
	routerA := edge.RouterA
	routerB := edge.RouterB

	// check that the routers are existing
	if routerA == nil || routerB == nil {
		log.Fatal().Msg("one of the edge nodes is nil")
		return
	}

	c.edgesLock.Lock()
	defer c.edgesLock.Unlock()

	// concatenate the bits
	edgeId := edge.EdgeId()
	c.edges[edgeId] = edge

	edge.Initialize()
}

func (c *OrionRegistryState) DispatchEdgeRemovedEvent(edge *Edge) {
	log.Debug().Msg("edge got removed")

	if c.edges[edge.EdgeId()] == nil {
		return
	}

	c.edgesLock.Lock()
	defer c.edgesLock.Unlock()

	// edge.RouterA.DispatchEdgeRemovedEvent(edge)
	// edge.RouterB.DispatchEdgeRemovedEvent(edge)

	c.edges[edge.EdgeId()].dispose()
	c.edges[edge.EdgeId()] = nil
}
