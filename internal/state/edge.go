package state

import (
	"context"
	"fmt"
	"slices"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EdgeIdentity uint64
type Edge struct {
	RouterA *Router
	RouterB *Router

	edgeObjectContext       context.Context
	edgeObjectContextCancel context.CancelCauseFunc

	globalState *OrionRegistryState
	log         zerolog.Logger
}

func NewEdge(
	globalContext context.Context,
	routerA *Router,
	routerB *Router,
	globalState *OrionRegistryState,
) *Edge {
	ctx, cancel := context.WithCancelCause(context.Background())
	edge := &Edge{
		RouterA: routerA,
		RouterB: routerB,

		edgeObjectContext:       ctx,
		edgeObjectContextCancel: cancel,

		globalState: globalState,
	}
	identity := edge.EdgeId()

	logger := log.With().Uint32("edge-identity", uint32(identity)).Logger()
	logger.Debug().Msg("new router session initiated")
	edge.log = logger

	return edge
}

func (c *Edge) EdgeId() EdgeIdentity {
	ids := []uint32{uint32(c.RouterA.Identity), uint32(c.RouterB.Identity)}
	slices.Sort(ids)
	return EdgeIdentity(uint64(ids[0]) << 32 & uint64(ids[1]))
}

// Sends a new initialization step to both the peers
// asking for a random one to choose a pre-shared key
// The hole-punching logic is done locally by the peers
func (c *Edge) Initialize() {
	go func() {
		select {
		case <-c.RouterA.routerObjectContext.Done():
		case <-c.RouterB.routerObjectContext.Done():
		case <-c.edgeObjectContext.Done():
		}

		// the edge lifetime is finished
		// todo: teardown edge
		c.globalState.DispatchEdgeRemovedEvent(c)
	}()

}

func (c *Edge) Dispose() {
	c.globalState.DispatchEdgeRemovedEvent(c)
}
func (c *Edge) dispose() {
	c.edgeObjectContextCancel(fmt.Errorf("edge is finished"))
}
