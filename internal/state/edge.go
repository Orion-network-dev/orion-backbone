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
	seeded                  chan struct{}

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
	logger.Debug().Msg("new edge session initiated")
	edge.log = logger

	return edge
}

func (c *Edge) EdgeId() EdgeIdentity {
	ids := []uint32{uint32(c.RouterA.Identity), uint32(c.RouterB.Identity)}
	slices.Sort(ids)
	return EdgeIdentity((uint64(ids[0])+1)<<32 + (uint64(ids[1]) + 1))
}

func (c *Edge) waitForSeed(token string) (*struct{}, error) {
	for {
		select {
		case <-c.RouterA.routerObjectContext.Done():
		case <-c.RouterB.routerObjectContext.Done():
		case <-c.edgeObjectContext.Done():
			c.seeded = nil
			c.log.Debug().Err(c.edgeObjectContext.Err()).Msg("edge context is finished")
			return nil, c.edgeObjectContext.Err()
		case data := <-c.seeded:
			c.seeded = nil
			return &data, nil
		}
	}

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
			c.log.Debug().Err(c.edgeObjectContext.Err()).Msg("edge context is finished")
			return
		}
		c.Dispose()
		c.log.Debug().Msg("starting edge disposal")
	}()

	c.log.Debug().Msg("edge instance started")

	c.seeded = make(chan struct{})
	// send a message to the routerA, requesting a new tunnel connection.
	c.RouterA.Dispatch(CreateEdgeRequest{})
	if _, err := c.waitForSeed("token!"); err == nil {

	} else {
		return
	}

	c.seeded = make(chan struct{})
	// send a message to the routerB, requesting a new tunnel connection.
	c.RouterB.Dispatch(CreateEdgeRequest{})
	if _, err := c.waitForSeed("token!"); err == nil {

	} else {
		return
	}

	c.seeded = make(chan struct{})
	c.RouterA.Dispatch(SeedEdgeRequest{})

	c.seeded = make(chan struct{})
	c.RouterB.Dispatch(SeedEdgeRequest{})
}

func (c *Edge) Dispose() {
	c.globalState.DispatchEdgeRemovedEvent(c)
}
func (c *Edge) dispose() {
	c.edgeObjectContextCancel(fmt.Errorf("edge is finished"))
}
