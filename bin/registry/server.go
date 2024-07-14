package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"google.golang.org/grpc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	debug = flag.Bool("debug", false, "change the log level to debug")
	port  = flag.Int("port", 6443, "the port the server will listen on")
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
		log.Fatal().Err(err).Msgf("Failed to read the required certificates")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to start listener")
	}

	// Create a new gRPC server
	s := grpc.NewServer(grpc.Creds(cred))

	registry, err := internal.NewOrionRegistryImplementation()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create the registry")
	}
	holepunch, err := internal.NewOrionHolePunchingImplementation()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create the holepunching service")
	}

	proto.RegisterRegistryServer(s, registry)
	proto.RegisterHolePunchingServiceServer(s, holepunch)

	// Start the gRPC server
	log.Info().Str("listening-address", lis.Addr().String()).Msgf("Server listening")
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}
