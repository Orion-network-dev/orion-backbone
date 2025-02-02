package server

import (
	"io/fs"
	"net/http"
)

type Server struct{}

func (c *Server) Handler() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve the static files
	fSys, _ := fs.Sub(fs.FS(assets), "static")
	fs := http.FileServer(http.FS(fSys))
	mux.Handle("/", fs)
	mux.HandleFunc("/whoami", c.whoami)
	mux.HandleFunc("/ws", c.upgrade)
	mux.HandleFunc("/state", c.state)

	return mux
}

func NewServer() *Server {
	return &Server{}
}
