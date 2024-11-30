package protocol

import (
	"context"
	"fmt"

	"github.com/MatthieuCoder/OrionV3/bin/registry-ws/server/protocol/messages"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/state"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var orionRegistryState *state.OrionRegistryState = state.NewOrionRegistryState()

type Client struct {
	ctx      context.Context
	router   *state.Router
	identity state.RouterIdentity
	ws       *websocket.Conn

	log zerolog.Logger
}

func NewClient(ws *websocket.Conn, identity state.RouterIdentity) *Client {
	c := &Client{
		ws:       ws,
		identity: identity,
		ctx:      context.Background(),
		log:      log.With().Uint32("router-identity", uint32(identity)).Logger(),
	}

	c.log.Debug().Msg("starting new client connection")
	go c.startRoutine()
	return c
}

func (c *Client) send(k string, msg state.Event) error {
	err := c.ws.WriteJSON(messages.Event{
		Kind:    k,
		Content: msg,
	})
	if err != nil {
		c.log.Error().Err(err).Msg("failed to send message")
	}
	return err
}

func (c *Client) startRoutine() {
	c.log.Debug().Msg("connection handling routine started")
	c.send(messages.MessageKindHello, messages.Hello{
		Message:  "Hi. This is orion-registry.",
		Identity: c.identity,
		Version:  internal.Version,
		Commit:   internal.Commit,
	})

	// check if the router exists
	rtrs := orionRegistryState.GetRouters()
	rtr := rtrs[c.identity]
	if rtr == nil {
		c.log.Debug().Msg("initialized a new state plane router object")
		rtr = state.NewRouter(context.Background(), c.identity, orionRegistryState)
		// dispatch new router if the given router doesn't exist
		orionRegistryState.DispatchNewRouterEvent(
			rtr,
		)
		c.router = rtr
	} else {
		c.log.Debug().Msg("re-using an existing router object system")
		// we inform the router object, in the registry state
		// that the connection is still ongoing and should not
		// be idle-disposed.
		c.router = rtr
	}

	sub := c.router.Subscribe()
	channel := sub.Ch()
	ctx, cancel := context.WithCancelCause(c.ctx)
	defer sub.Close()

	go func() {
		// listening for events
		for {
			select {
			case event := <-channel:
				switch event := event.(type) {
				case state.RouterConnect:
					c.log.Debug().Msg("sending a new router connect event")
					c.send(messages.MessageKindRouterConnect, event)
				case state.RouterDisconnect:
					c.log.Debug().Msg("sending a new disconnect event")
					c.send(messages.MessageKindRouterDisconnect, event)
				}
			case <-c.ctx.Done():
				goto end
			}
		}

	end:
		cancel(fmt.Errorf("the state-internal event listener is finished"))
	}()

	// We start listening for events once the listener go-routine is setup
	// this is because the increment connection count trigers a recovery
	// of a lost connection
	rtr.IncrementConnectionCount()

	go func() {
		for {
			_, _, err := c.ws.ReadMessage()
			if err != nil {
				goto end
			}

			// todo: handle clients events
		}

	end:
		cancel(fmt.Errorf("the websocket listening is finished"))
	}()

	// wait for the context to be finished
	<-ctx.Done()

	c.log.Info().Msg("connection ended")
	c.router.DecrementConnectionCount()

	c.ws.Close()
}
