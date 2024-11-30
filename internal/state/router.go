package state

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/teivah/broadcast"
)

type RouterIdentity uint32
type Router struct {
	Identity RouterIdentity

	sending *broadcast.Relay[Event]

	routerObjectContext            context.Context
	routerObjectContextCancel      context.CancelCauseFunc
	connectionsCount               atomic.Int32
	connectionTimeoutContextCancel context.CancelCauseFunc

	globalState *OrionRegistryState
	log         zerolog.Logger
}

func NewRouter(
	globalContext context.Context,
	identity RouterIdentity,
	globalState *OrionRegistryState,
) *Router {
	ctx, cancel := context.WithCancelCause(globalContext)
	return &Router{
		Identity:                  identity,
		connectionsCount:          atomic.Int32{},
		sending:                   broadcast.NewRelay[Event](),
		routerObjectContext:       ctx,
		routerObjectContextCancel: cancel,
		globalState:               globalState,
		log:                       log.With().Uint32("router-identity", uint32(identity)).Logger(),
	}
}

func (c *Router) Subscribe() *broadcast.Listener[Event] {
	return c.sending.Listener(1)
}

func (c *Router) IncrementConnectionCount() {
	c.connectionsCount.Add(1)
	c.updateConnectionsCountRoutine()
}
func (c *Router) DecrementConnectionCount() {
	c.connectionsCount.Add(-1)
	c.updateConnectionsCountRoutine()
}

func (c *Router) updateConnectionsCountRoutine() {
	current := c.connectionsCount.Load()

	c.log.Debug().
		Int32("connections-count", current).
		Msg("connection count ended, updating")

	if c.connectionTimeoutContextCancel != nil {
		c.connectionTimeoutContextCancel(fmt.Errorf("connections count updated, canceling previous tasks"))
	}

	if current == 0 {
		ctx, e := context.WithCancelCause(context.Background())
		c.connectionTimeoutContextCancel = e

		// TODO: we must implement a way to deal with not-sent messages

		// we launch a background task ticking for either
		// 	1. the session timeout mechanism
		//	2. a new connection
		go func() {
			timeout := time.NewTimer(time.Minute)
			c.log.Debug().Msg("ticking a minute before session expiration")

			subscribe := c.Subscribe()
			replayPending := make([]*Event, 1000)
			replayEventsCount := 0

			for {
				select {
				case message := <-subscribe.Ch():
					log.Debug().Msg("appending to replay pending messages")
					replayPending[replayEventsCount] = &message
					replayEventsCount += 1
				case <-ctx.Done():
					subscribe.Close()

					c.log.Debug().
						Int("events-count", replayEventsCount).
						Msg("session resumed, replaying events")

					// dequeue events
					for i, event := range replayPending {
						c.sending.Broadcast(event)
						if i == replayEventsCount {
							break
						}
					}
					goto end
				case <-timeout.C:
					log.Debug().Msg("session expired")
					subscribe.Close()
					c.globalState.DispatchRouterRemovedEvent(c)
					goto end
				}
			}
		end:
			c.connectionTimeoutContextCancel = nil
		}()
	}

}

func (c *Router) DispatchNewRouterEvent(router *Router) {
	newRouterEvent := RouterConnect{
		Router: router,
	}
	c.sending.Broadcast(newRouterEvent)
}

func (c *Router) Dispose() {
	log.Info().Msg("disposing of router")
}
