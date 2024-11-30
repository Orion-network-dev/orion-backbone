package state

import (
	"sync"
)

type OrionRegistryRoutersState map[RouterIdentity]*Router
type OrionRegistryEdgesState map[uint64]*Edge

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
