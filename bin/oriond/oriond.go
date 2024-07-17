package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var (
	debug          = flag.Bool("debug", false, "change the log level to debug")
	registryServer = flag.String("registry-server", "reg.orionet.re", "the address of the registry server")
	registryPort   = flag.Uint("registry-port", 443, "the port used by the registry")
)

func main() {
	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Get TLS credentials
	cred, err := internal.LoadTLS(false)
	if err != nil {
		log.Error().Err(err).Msgf("unable to connect gRPC channel")
		return
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", *registryServer, *registryPort),
		grpc.WithTransportCredentials(cred),
		grpc.WithIdleTimeout(time.Second*120),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to start the grpc client")
		return
	}

	_, err = implementation.NewOrionClientDaemon(
		context.Background(),
		conn,
	)

	if err != nil {
		log.Error().Err(err).Msgf("failed to bring up orion client daemon")
		return
	}
}
