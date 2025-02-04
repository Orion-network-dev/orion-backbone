package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/orion-network-dev/orion-backbone/bin/oriond/implementation"
	"github.com/orion-network-dev/orion-backbone/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	debug          = flag.Bool("debug", false, "change the log level to debug")
	registryServer = flag.String("registry-server", "reg.orionet.re:64431", "the address of the registry server")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// listen for interrupts
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// parse all command flags
	flag.Parse()

	// setup the time format used by the logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// if the debug lag is used, we set the logging level to debug
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// we print the version information
	internal.PrintVersionHeader()

	oriond, err := implementation.NewOrionClientDaemon(ctx)
	defer oriond.Dispose()

	// We load the required certificates
	privateKey, chain := internal.LoadPemFile()
	certificateKeyPair := internal.LoadX509KeyPair(privateKey, chain)
	authorityPool, err := internal.LoadAuthorityPool()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load the required certificates")
		return
	}

	// information required to connect to the registry over websocket
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{certificateKeyPair},
			RootCAs:      authorityPool,
			MinVersion:   tls.VersionTLS13,
			MaxVersion:   tls.VersionTLS13,
		},
		Subprotocols: []string{fmt.Sprintf("orion-registry-%s", internal.Commit)},
	}
	url := url.URL{
		Scheme: "wss",
		Host:   *registryServer,
		Path:   "/ws",
	}

	// we dial the registry server, initializing the tls1.3 connection
	connection, _, err := dialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("unsable to dial")
	}
	defer connection.Close()

	// we initialize the oriond daemon to handle the websocket messages
	oriond.ListenOnWS(connection)

	select {
	case <-ctx.Done():
		return
	case <-interrupt:
		cancel()
		<-ctx.Done()
		return
	}
}
