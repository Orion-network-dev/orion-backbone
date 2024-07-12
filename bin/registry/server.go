package main

import (
	"flag"
	"net"
	"os"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"google.golang.org/grpc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	cred := internal.NewServerTLS()

	// Create a listener that listens to localhost port 8443
	lis, err := net.Listen("tcp", ":6443")

	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to start listener")
	}

	// Create a new gRPC server
	s := grpc.NewServer(grpc.Creds(cred))

	// Register the Registry service.
	proto.RegisterRegistryServer(s, &internal.OrionRegistryImplementations{})

	// Create the hole-punching service
	holepunch, err := internal.NewOrionHolePunchingImplementations()
	if err != nil {
		panic(err)
	}
	// Register the hole-punching service
	proto.RegisterHolePunchingServiceServer(s, holepunch)

	// Start the gRPC server
	log.Info().Str("listening-address", lis.Addr().String()).Msgf("Server listening")
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}
