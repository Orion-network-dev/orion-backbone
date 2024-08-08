package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MatthieuCoder/OrionV3/bin/oriond/implementation"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	debug          = flag.Bool("debug", false, "change the log level to debug")
	registryServer = flag.String("registry-server", "reg.orionet.re", "the address of the registry server")
	pprof          = flag.String("debug-pprof", "0.0.0.0:6060", "")
	registryPort   = flag.Uint("registry-port", 6443, "the port used by the registry")
)

func main() {
	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		go func() {
			fmt.Println(http.ListenAndServe(*pprof, nil))
		}()
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
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 20,
			Timeout:             time.Second * 1,
			PermitWithoutStream: false,
		}),
	)
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
	err = orionDaemon.Start()
	if err != nil {
		log.Error().Err(err).Msgf("failed to bring up orion client daemon")
		return
	}

	defer orionDaemon.Dispose()
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	restartCounter := 0
	for {
		select {
		case <-sigs:
			return
		case <-ctx.Done():
			return
		case <-orionDaemon.Context.Done():
			restartCounter += 1
			if restartCounter > 10 {
				return
			}
			time.Sleep(time.Second * 5)
			err := orionDaemon.Start()
			if err != nil {
				log.Error().Err(err).Msgf("failed to bring up orion client daemon")
				return
			}
		}
	}
}
