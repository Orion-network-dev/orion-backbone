package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/MatthieuCoder/OrionV3/bin/registry-ws/server"
	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	tlsConfig := &tls.Config{
		Certificates:             []tls.Certificate{certificateKeyPair},
		ClientCAs:                authorityPool,
		ClientAuth:               tls.VerifyClientCertIfGiven,
		MinVersion:               tls.VersionTLS13,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}

	srv := server.NewServer()

	server := http.Server{
		Addr:      *listeningHost,
		Handler:   srv.Handler(),
		TLSConfig: tlsConfig,
	}

	ln, err := net.Listen("tcp", *listeningHost)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	tlsListener := tls.NewListener(ln, tlsConfig)
	if err := server.Serve(tlsListener); err != nil {
		log.Fatal().Msgf("(HTTPS) error listening to port: %v", err)
	}
}
