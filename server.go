package gbio

import (
	"net/http"

	"github.com/anqur/gbio/internal/servers"
)

type Server struct {
	s *servers.Server
}

type ServerOption func(s *servers.Server)

func NewServer(addr string, opts ...ServerOption) *Server {
	s := &servers.Server{Server: http.Server{Addr: addr}}
	for _, opt := range opts {
		opt(s)
	}
	return &Server{s: s}
}

var DefaultServer = NewServer(":8080")

func WithMux(m http.Handler) ServerOption {
	return func(s *servers.Server) { s.Handler = m }
}
func UseMux(m http.Handler) { DefaultServer.s.Handler = m }

func UseAddr(addr string) { DefaultServer.s.Addr = addr }

func ListenAndServe() error { return DefaultServer.s.ListenAndServe() }
