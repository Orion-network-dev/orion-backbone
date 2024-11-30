package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/MatthieuCoder/OrionV3/bin/registry-ws/server/protocol"
	"github.com/MatthieuCoder/OrionV3/internal/state"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{}

func upgradeErrorPage(w http.ResponseWriter) {
	file, _ := errors.ReadFile("auth-required.html")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/html")
	w.Write(file)
}

func (c *Server) upgrade(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		upgradeErrorPage(w)
		return
	}

	cz, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade a http(s) connection to a websocket connection")
		return
	}
	// todo: close websocket connections cleanly
	//	defer cz.Close()

	leaf := r.TLS.PeerCertificates[0]
	cn := leaf.Subject.CommonName
	cnParts := strings.Split(cn, ":")
	if len(cnParts) != 2 || cnParts[1] != "oriond" {
		log.Error().Err(err).Msg("The given certificate is not valid for logging-in into oriond")
		return
	}

	routerId, err := strconv.Atoi(cnParts[0])
	if err != nil {
		log.Error().Err(err).Msg("The given certificate is not valid for logging-in into oriond")
		return
	}
	identity := state.RouterIdentity(routerId)

	protocol.NewClient(cz, identity)
}
