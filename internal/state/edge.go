package state

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"slices"
	"time"

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
	logger.Debug().Msg("new edge session initiated")
	edge.log = logger

	return edge
}

func (c *Edge) EdgeId() EdgeIdentity {
	ids := []uint32{uint32(c.RouterA.Identity), uint32(c.RouterB.Identity)}
	slices.Sort(ids)
	return EdgeIdentity((uint64(ids[0])+1)<<32 + (uint64(ids[1]) + 1))
}

func (c *Edge) seedEdgeAndWait(
	router *Router,
	othrIdentity RouterIdentity,
) (*CreateEdgeResponse, error) {
	chanCallback := make(chan CreateEdgeResponse)
	router.edgeResponseCallback = &chanCallback
	// send a message to the routerA, requesting a new tunnel connection.
	router.Dispatch(CreateEdgeRequest{
		PeerID: othrIdentity,
		EdgeID: c.EdgeId(),
	})

	ticker := time.NewTimer(time.Second * 60)
	select {
	case <-c.RouterA.routerObjectContext.Done():
		goto end
	case <-c.RouterB.routerObjectContext.Done():
		goto end
	case <-ticker.C:
		goto end
	case <-c.edgeObjectContext.Done():
		goto end
	case data := <-chanCallback:
		router.edgeResponseCallback = nil
		return &data, nil
	}

end:
	router.edgeResponseCallback = nil
	err := fmt.Errorf("timeout reached while waiting for seeding")
	c.log.Debug().Err(err).Msg("timeout reached while waiting")
	return nil, err
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

	seedRouterA, err := c.seedEdgeAndWait(c.RouterA, c.RouterB.Identity)
	if err != nil {
		c.log.Err(err)
		c.Dispose()
		return
	}
	c.log.Debug().Msg("edge seeding (1/2)")

	seedRouterB, err := c.seedEdgeAndWait(c.RouterB, c.RouterA.Identity)
	if err != nil {
		c.log.Err(err)
		c.Dispose()
		return
	}

	c.log.Debug().Msg("got initialzation data, seeding edges")

	bytesA, err := base64.StdEncoding.DecodeString(seedRouterA.PresharedKeybB4)
	if err != nil {
		c.log.Err(err)
		c.Dispose()
		return
	}
	bytesB, err := base64.StdEncoding.DecodeString(seedRouterB.PresharedKeybB4)
	if err != nil {
		c.log.Err(err)
		c.Dispose()
		return
	}

	var random [32]byte
	for i := range 32 {
		token := make([]byte, 1)
		rand.Read(token)
		random[i] = bytesA[i] ^ bytesB[i] ^ token[0]
	}

	presharedKey := base64.StdEncoding.EncodeToString(random[:])

	c.log.Debug().Msg("seeding edges")
	c.RouterA.Dispatch(SeedEdgeRequest{
		OtherPeer:    seedRouterB.PublicEndpoint,
		PresharedKey: presharedKey,
	})
	c.RouterB.Dispatch(SeedEdgeRequest{
		OtherPeer:    seedRouterB.PublicEndpoint,
		PresharedKey: presharedKey,
	})
}

func (c *Edge) Dispose() {
	c.log.Debug().Msg("dispose")
	c.globalState.DispatchEdgeRemovedEvent(c)
}
func (c *Edge) dispose() {
	c.log.Debug().Msg("internal: disposing")
	c.edgeObjectContextCancel(fmt.Errorf("edge is finished"))
}
