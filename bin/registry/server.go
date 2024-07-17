package main

import (
	"flag"
	"net"
	"os"

	"github.com/MatthieuCoder/OrionV3/bin/registry/implementation"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var (
	debug         = flag.Bool("debug", false, "change the log level to debug")
	listeningHost = flag.String("listen-host", "127.0.0.1:6443", "the port the server will listen on")
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

	cred, err := internal.LoadTLS(true)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read the required certificates")
		return
	}

	lis, err := net.Listen("tcp", *listeningHost)

	if err != nil {
		log.Error().Err(err).Msgf("Failed to start listener")
		return
	}

	// Create a new gRPC server
	s := grpc.NewServer(grpc.Creds(cred))

	registry, err := implementation.NewOrionRegistryImplementation()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create the registry")
		return
	}
	holepunch, err := implementation.NewOrionHolePunchingImplementation()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create the holepunching service")
		return
	}

	proto.RegisterRegistryServer(s, registry)
	proto.RegisterHolePunchingServiceServer(s, holepunch)

	// Start the gRPC server
	log.Info().Str("listening-address", lis.Addr().String()).Msgf("Server listening")
	if err := s.Serve(lis); err != nil {
		log.Error().Err(err).Msg("Failed to serve")
		return
	}
}
