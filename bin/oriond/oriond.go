package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	debug          = flag.Bool("debug", false, "change the log level to debug")
	registryServer = flag.String("registry-server", "reg.orionet.re", "the address of the registry server")
	registryPort   = flag.Uint("registry-port", 443, "the port used by the registry")
)

func main() {
	flag.Parse()

	cred, err := internal.LoadTLS(false)
	if err != nil {
		log.Error().Err(err).Msgf("unable to connect gRPC channel")
		return
	}

	conn, err := createGRPCClientConnection(cred)
	if err != nil {
		log.Error().Err(err).Msg("failed to start the grpc client")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	orionDaemon, err := implementation.NewOrionClientDaemon(
		ctx,
		conn,
	)
	if err != nil {
		log.Error().Err(err).Msgf("failed to bring up orion client daemon")
		return
	}
	defer orionDaemon.Dispose()
	defer cancel()

	// Wait for exit signal
	waitForExitSignal()
}

func waitForExitSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func setupLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if *debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func createGRPCClientConnection(cred credentials.TransportCredentials) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		fmt.Sprintf("%s:%d", *registryServer, *registryPort),
		grpc.WithTransportCredentials(cred),
		grpc.WithIdleTimeout(time.Second*120),
	)
}
