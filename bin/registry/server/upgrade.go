package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/orion-network-dev/orion-backbone/bin/registry/server/protocol"
	"github.com/orion-network-dev/orion-backbone/internal/state"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{}

func upgradeErrorPage(w http.ResponseWriter) {
	file, _ := errors.ReadFile("auth-required.html")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/html")
	w.Write(file)
}

func (c *Server) state(w http.ResponseWriter, r *http.Request) {
	s, _ := json.Marshal(protocol.OrionRegistryState)
	w.Header().Add("Content-Type", "application/json")
	w.Write(s)
}

func (c *Server) upgrade(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		upgradeErrorPage(w)
		return
	}

	cz, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to upgrade a http(s) connection to a websocket connection")
		return
	}

	leaf := r.TLS.PeerCertificates[0]
	cn := leaf.Subject.CommonName
	cnParts := strings.Split(cn, ":")
	if len(cnParts) != 2 || cnParts[1] != "oriond" {
		log.Error().Err(err).Msg("the given certificate is not valid for logging-in into oriond")
		return
	}

	routerId, err := strconv.Atoi(cnParts[0])
	if err != nil {
		log.Error().Err(err).Msg("the given certificate is not valid for logging-in into oriond")
		return
	}
	identity := state.RouterIdentity(routerId)

	sessionId := r.Header.Get("X-SessionId")

	protocol.NewClient(cz, identity, sessionId)
}
