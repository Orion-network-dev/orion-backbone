package main

import (
	"flag"
	"net"
	"os"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/getsentry/sentry-go"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_sentry "github.com/johnbellone/grpc-middleware-sentry"
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

	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://5c2733c29efd992fbcbd9d01b3b0ab8e@o228322.ingest.us.sentry.io/4507616547569664",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatal().Msgf("sentry.Init: %s", err)
	}

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

	lis, err := net.Listen("tcp", *listeningHost)

	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to start listener")
	}

	// Create a new gRPC server
	s := grpc.NewServer(grpc.Creds(cred),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_sentry.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_sentry.UnaryServerInterceptor(),
		)))

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
