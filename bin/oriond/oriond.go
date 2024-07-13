package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")
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
		log.Fatal().Msgf("Unable to connect gRPC channel %v", err)
	}

	_, err = grpc.NewClient(fmt.Sprintf("%s:%d", "reg.orionet.re", 6443), grpc.WithTransportCredentials(cred), grpc.WithIdleTimeout(time.Second*120))
	if err != nil {
		log.Fatal().Msgf("Unable to connect gRPC channel %v", err)
	}

	// Create the gRPC client
	// registryClient := proto.NewRegistryClient(conn)
	// holepunchingClient = proto.NewHolePunchingServiceClient(conn)
}
