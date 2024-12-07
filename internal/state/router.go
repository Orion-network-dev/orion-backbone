package state

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/teivah/broadcast"
	"golang.org/x/exp/rand"
)

type RouterIdentity uint32
type Router struct {
	Identity RouterIdentity

	sending *broadcast.Relay[Event]

	routerObjectContext            context.Context
	routerObjectContextCancel      context.CancelCauseFunc
	connectionsCount               atomic.Int32
	connectionTimeoutContextCancel context.CancelCauseFunc

	session string

	globalState *OrionRegistryState
	log         zerolog.Logger
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NewRouter(
	globalContext context.Context,
	identity RouterIdentity,
	globalState *OrionRegistryState,
) *Router {
	ctx, cancel := context.WithCancelCause(globalContext)
	session := randStringBytes(128)
	logger := log.With().Uint32("router-identity", uint32(identity)).Logger()
	logger.Debug().Str("session", session).Msg("new router session initiated")
	return &Router{
		Identity:                  identity,
		connectionsCount:          atomic.Int32{},
		sending:                   broadcast.NewRelay[Event](),
		routerObjectContext:       ctx,
		routerObjectContextCancel: cancel,
		globalState:               globalState,
		session:                   session,
		log:                       logger,
	}
}

func (c *Router) SessionId() string {
	return c.session
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
		Msg("connection count ended, updating")

	if c.connectionTimeoutContextCancel != nil {
		c.connectionTimeoutContextCancel(fmt.Errorf("connections count updated, canceling previous tasks"))
	}

	if current == 0 {
		ctx, e := context.WithCancelCause(context.Background())
		c.connectionTimeoutContextCancel = e

		// we launch a background task ticking for either
		// 	1. the session timeout mechanism
		//	2. a new connection
		go func() {
			timeout := time.NewTimer(time.Minute)
			c.log.Debug().Msg("ticking a minute before session expiration")

			subscribe := c.Subscribe()
			defer subscribe.Close()

			replayPending := make([]Event, 1000)
			replayEventsCount := 0

			for {
				select {
				case message := <-subscribe.Ch():
					log.Debug().Msg("appending to replay pending messages")
					replayPending[replayEventsCount] = message
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
					c.Dispose()
					goto end
				case <-c.routerObjectContext.Done():
					goto end
				}
			}
		end:
			c.connectionTimeoutContextCancel = nil
		}()
	}

}

func (c *Router) DispatchNewRouterEvent(router *Router) {
	c.sending.Broadcast(RouterConnectEvent{
		Router: router,
	})
}

func (c *Router) DispatchRouterRemovedEvent(router *Router) {
	c.sending.Broadcast(RouterDisconnectEvent{
		Router: router,
	})
}

func (c *Router) DispatchNewEdgeEvent(edge *Edge) {

}
func (c *Router) DispatchEdgeRemovedEvent(edge *Edge) {}

func (c *Router) dispose() {
	c.routerObjectContextCancel(fmt.Errorf("router is disposed"))
	c.log.Debug().Msg("context canceled")
}

func (c *Router) Dispose() {
	c.globalState.DispatchRouterRemovedEvent(c)
}
