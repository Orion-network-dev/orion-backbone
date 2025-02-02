package protocol

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/state"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var OrionRegistryState *state.OrionRegistryState = state.NewOrionRegistryState()

type Client struct {
	ctx      context.Context
	router   *state.Router
	identity state.RouterIdentity
	ws       *websocket.Conn

	log zerolog.Logger
}

func NewClient(ws *websocket.Conn, identity state.RouterIdentity, sessionId string) *Client {
	c := &Client{
		ws:       ws,
		identity: identity,
		ctx:      context.Background(),
		log:      log.With().Uint32("router-identity", uint32(identity)).Logger(),
	}
	c.log.Debug().Msg("starting new client connection")
	go c.startRoutine(sessionId)
	return c
}

func (c *Client) send(event *state.JsonEvent) error {
	err := c.ws.WriteJSON(event)
	if err != nil {
		c.log.Error().Err(err).Msg("failed to send message")
	}
	return err
}

func (c *Client) startRoutine(sessionId string) {
	defer c.ws.Close()
	c.log.Debug().Msg("connection handling routine started")

	// check if the router exists
	rtrs := OrionRegistryState.GetRouters()
	rtr := rtrs[c.identity]
	if rtr == nil {
		c.log.Debug().Msg("initialized a new state plane router object")
		rtr = state.NewRouter(context.Background(), c.identity, OrionRegistryState)
		// dispatch new router if the given router doesn't exist
		OrionRegistryState.DispatchNewRouterEvent(
			rtr,
		)
	} else {
		if rtr.SessionId() == sessionId {
			c.log.Debug().Msg("re-using an existing router object system")
			// we inform the router object, in the registry state
			// that the connection is still ongoing and should not
			// be idle-disposed.
			c.router = rtr
		} else {
			c.log.Debug().Msg("deleted old session, initializing new session")
			rtr.Dispose()
			rtr = state.NewRouter(context.Background(), c.identity, OrionRegistryState)
			// dispatch new router if the given router doesn't exist
			OrionRegistryState.DispatchNewRouterEvent(
				rtr,
			)
		}
	}
	c.router = rtr
	event, _ := state.MarshalEvent(state.Hello{
		Message:  "Hi. This is orion-registry.",
		Identity: c.router.Identity,
		Version:  internal.Version,
		Commit:   internal.Commit,
		Session:  c.router.SessionId(),
	})
	// we send the hello message
	c.send(event)

	ctx, cancel := context.WithCancelCause(c.ctx)

	go func() {
		// subscribe to the client events from the state
		sub := c.router.Subscribe()
		channel := sub.Ch()
		defer sub.Close()

		// listening for events on the stream
		for {
			select {
			case event := <-channel:
				c.send(event)
			case <-ctx.Done():
				c.log.Debug().Msg("server state listening routine is done")
				return
			}
		}
	}()

	// We start listening for events once the listener go-routine is setup
	// this is because the increment connection count trigers a recovery
	// of a lost connection
	rtr.IncrementConnectionCount()
	defer c.router.DecrementConnectionCount()

	go func() {
		for {
			_, data, err := c.ws.ReadMessage()
			if err != nil {
				goto end
			}

			event := state.JsonEvent{}
			if err := json.Unmarshal(data, &event); err != nil {
				c.log.Error().Err(err).Msg("failed to parse jsonevent")
				goto end
			}
			out, err := state.UnmarshalEvent(event)
			if err != nil {
				c.log.Error().Err(err).Msg("failed to parse event")
				goto end
			}

			switch message := out.(type) {
			// sent once a router wants to connect to another one
			case state.RouterInitiateRequest:
				c.log.Info().Msgf("received a router connect event to %d", *message.Identity)

				routers := OrionRegistryState.GetRouters()
				if routers[*message.Identity] == nil || routers[c.identity] == nil {
					goto end
				}

				OrionRegistryState.DispatchNewEdge(
					state.NewEdge(context.Background(), routers[*message.Identity], routers[c.identity], OrionRegistryState),
				)

				continue
			default:
				c.log.Error().Str("event", event.Kind).Msg("unknown event type")
				goto end
			}
		}

	end:
		cancel(fmt.Errorf("the websocket listening is finished"))
	}()

	// wait for the context to be finished
	<-ctx.Done()

	c.log.Info().Err(context.Cause(ctx)).Msg("connection routine ended")
}
