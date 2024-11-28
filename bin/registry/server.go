package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/MatthieuCoder/OrionV3/bin/registry/implementation"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	pprof         = flag.String("debug-pprof", ":6060", "")
	enable_prof   = flag.Bool("enable-pprof", false, "enable pprof for debugging")
	debug         = flag.Bool("debug", false, "change the log level to debug")
	listeningHost = flag.String("listen-host", "127.0.0.1:6443", "the port the server will listen on")
)

func main() {
	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	flag.Parse()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *enable_prof {
		go func() {
			fmt.Println(http.ListenAndServe(*pprof, nil))
		}()
	}
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	privateKey, chain := internal.LoadPemFile()
	certificateKeyPair := internal.LoadX509KeyPair(privateKey, chain)
	authorityPool, err := internal.LoadAuthorityPool()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to start listener")
		return
	}

	keyChain := credentials.NewTLS(
		&tls.Config{
			Certificates: []tls.Certificate{certificateKeyPair},
			RootCAs:      authorityPool,
			MinVersion:   tls.VersionTLS13,
			MaxVersion:   tls.VersionTLS13,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    authorityPool,
		},
	)

	lis, err := net.Listen("tcp", *listeningHost)

	if err != nil {
		log.Error().Err(err).Msgf("Failed to start listener")
		return
	}

	// Create a new gRPC server
	s := grpc.NewServer(
		grpc.Creds(keyChain),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:              time.Second * 20,
			Timeout:           time.Second * 1,
			MaxConnectionIdle: time.Second * 20,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             time.Second * 15,
			PermitWithoutStream: false,
		}),
	)

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
