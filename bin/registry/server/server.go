package server

import (
	"io/fs"
	"net/http"

	"github.com/orion-network-dev/orion-backbone/internal"
)

type Server struct {
	tasksAssigner internal.LockableTasks
}

func (c *Server) Handler() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve the static files
	fSys, _ := fs.Sub(fs.FS(assets), "static")
	fs := http.FileServer(http.FS(fSys))
	mux.Handle("/", fs)
	mux.HandleFunc("/whoami", c.whoami)
	mux.HandleFunc("/ws", c.upgrade)
	mux.HandleFunc("/state", c.state)
	mux.HandleFunc("/holepunch", c.upgradeHolepunch)

	return mux
}

func NewServer() *Server {
	return &Server{}
}
