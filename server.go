package gbio

import (
	"net/http"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/anqur/gbio/internal/registries"
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

func WithServiceRegistry(
	cfg *etcd.Config,
	opts ...ServiceRegistryOption,
) ServerOption {
	return func(s *servers.Server) {
		s.Reg = registries.NewRegistry(cfg)
		for _, opt := range opts {
			opt(s.Reg)
		}
	}
}
func UseServiceRegistry(c *etcd.Config, opts ...ServiceRegistryOption) {
	DefaultServer.s.Reg = registries.NewRegistry(c)
	for _, opt := range opts {
		opt(DefaultServer.s.Reg)
	}
}

func UseAddr(addr string) { DefaultServer.s.Addr = addr }

func ListenAndServe() error { return DefaultServer.s.ListenAndServe() }
