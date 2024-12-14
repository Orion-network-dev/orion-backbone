package server

import (
	"net/http"
)

type Server struct{}

func (c *Server) Handler() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve the static files
	fs := http.FileServer(http.FS(assets))
	mux.Handle("/", http.StripPrefix("/static", fs))
	mux.HandleFunc("/whoami", c.whoami)
	mux.HandleFunc("/ws", c.upgrade)

	return mux
}

func NewServer() *Server {
	return &Server{}
}
