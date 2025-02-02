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
	Routers     OrionRegistryRoutersState `json:"routers"` // List of routers
	routersLock sync.Mutex

	Edges     OrionRegistryEdgesState `json:"edges"` // List of edges in the orion graph
	edgesLock sync.Mutex
}

func NewOrionRegistryState() *OrionRegistryState {
	return &OrionRegistryState{
		Routers:     make(OrionRegistryRoutersState),
		routersLock: sync.Mutex{},
		Edges:       make(OrionRegistryEdgesState),
		edgesLock:   sync.Mutex{},
	}
}

func (c *OrionRegistryState) GetRouters() OrionRegistryRoutersState {
	return c.Routers
}

// Called once a new member is joining the network
func (c *OrionRegistryState) DispatchNewRouterEvent(
	newRouter *Router,
) {
	c.routersLock.Lock()
	defer c.routersLock.Unlock()

	c.Routers[RouterIdentity(newRouter.Identity)] = newRouter

	for _, router := range c.Routers {
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
	if c.Routers[deletedRouter.Identity] == nil {
		return
	}

	c.routersLock.Lock()
	defer c.routersLock.Unlock()

	c.Routers[deletedRouter.Identity].dispose()
	delete(c.Routers, deletedRouter.Identity)
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
	edgeId := edge.EdgeId()
	c.Edges[edgeId] = edge
	c.edgesLock.Unlock()

	go edge.Initialize()
}

func (c *OrionRegistryState) DispatchEdgeRemovedEvent(edge *Edge) {
	if c.Edges[edge.EdgeId()] == nil {
		return
	}
	log.Debug().Int32("edge-id", int32(edge.EdgeId())).Msg("edge got removed")

	c.edgesLock.Lock()
	defer c.edgesLock.Unlock()

	edgeRemovedEvent := RouterEdgeRemovedEvent{
		Edge: edge,
	}
	edge.RouterA.Dispatch(edgeRemovedEvent)
	edge.RouterB.Dispatch(edgeRemovedEvent)

	c.Edges[edge.EdgeId()].dispose()
	delete(c.Edges, edge.EdgeId())
}
