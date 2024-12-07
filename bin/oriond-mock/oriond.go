package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/registry/server/protocol/messages"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/state"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	enable_prof    = flag.Bool("enable-pprof", false, "enable pprof for debugging")
	debug          = flag.Bool("debug", false, "change the log level to debug")
	registryServer = flag.String("registry-server", "reg.orionet.re", "the address of the registry server")
	pprof          = flag.String("debug-pprof", "0.0.0.0:6061", "")
	registryPort   = flag.Uint("registry-port", 6443, "the port used by the registry")
)

type Event struct {
	Kind    string          `json:"k"`
	Content json.RawMessage `json:"e"`
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *enable_prof {
		go func() {
			fmt.Println(http.ListenAndServe(*pprof, nil))
		}()
	}
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	url := url.URL{
		Scheme: "wss",
		Host:   "reg.orionet.re:6443",
		Path:   "/ws",
	}

	privateKey, chain := internal.LoadPemFile()
	certificateKeyPair := internal.LoadX509KeyPair(privateKey, chain)
	authorityPool, err := internal.LoadAuthorityPool()
	if err != nil {
		panic(err)
	}

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{certificateKeyPair},
			RootCAs:      authorityPool,
			MinVersion:   tls.VersionTLS13,
			MaxVersion:   tls.VersionTLS13,
		},
		Subprotocols: []string{"orion-reg-rpc"},
	}

	c, _, err := dialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal().Msgf("dial: %s", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Print("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			msg := Event{}
			json.Unmarshal(message, &msg)

			log.Printf("received %s... handling", msg.Kind)

			switch msg.Kind {
			case messages.MessageKindHello:
				hello := messages.Hello{}
				if err := json.Unmarshal(msg.Content, &hello); err != nil {
					panic("invalid json")
				}
				log.Printf("Hello message: %s", hello.Message)
				continue
			case messages.MessageKindRouterConnect:
				router_connect := state.RouterConnectEvent{}
				if err := json.Unmarshal(msg.Content, &router_connect); err != nil {
					panic("invalid json")
				}
				log.Printf("Router connected: %d", router_connect.Router.Identity)

				// todo: send request after provisionning the wireguard request

				continue
			case messages.MessageKindRouterDisconnect:
				router_connect := state.RouterDisconnectEvent{}
				if err := json.Unmarshal(msg.Content, &router_connect); err != nil {
					panic("invalid json")
				}
				log.Printf("Router disconnected: %d", router_connect.Router.Identity)
				continue
			case messages.MessageKindRouterEdgeInitConnectRequest:
				log.Printf("connect request received")
			case messages.MessageKindRouterEdgeInitConnectRequestResponse:
				log.Printf("connect request ack received")
				continue
			default:
				log.Printf("unroceverable error")
				panic("invalid kind")
			}

		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

		case <-done:
			return
		case <-interrupt:
			log.Print("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Print("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}
