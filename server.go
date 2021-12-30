package gbio

import (
	"fmt"
	"net/http"
	"net/url"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/anqur/gbio/internal/errors"
	"github.com/anqur/gbio/internal/registries"
	"github.com/anqur/gbio/internal/servers"
	"github.com/anqur/gbio/internal/utils"
)

var (
	ErrServerNotService = fmt.Errorf("%w: not a service, you might forget the codegen", errors.Err)
)

type Server struct {
	u *url.URL
	s *servers.Server
}

type ServerOption func(s *servers.Server)

func NewServer(host string, opts ...ServerOption) (*Server, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	s := &servers.Server{Server: http.Server{Addr: u.Host}}
	for _, opt := range opts {
		opt(s)
	}
	return &Server{u: u, s: s}, nil
}

var DefaultServer, _ = NewServer("localhost:8080")

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

func ListenAndServe() error { return DefaultServer.ListenAndServe() }

func (s *Server) ListenAndServe() error {
	if s.s.Reg != nil {
		if err := s.Register(); err != nil {
			return err
		}
	}
	return s.s.ListenAndServe()
}

func (s *Server) Addr() (string, error) {
	ip, err := utils.GetIPv4()
	if err != nil {
		return "", err
	}
	port := s.u.Port()
	if port == "" {
		return ip, nil
	}
	return fmt.Sprintf("%s:%s", ip, port), nil
}

func (s *Server) Register() error {
	srv, ok := s.s.Handler.(servers.Service)
	if !ok {
		return ErrServerNotService
	}
	addr, err := s.Addr()
	if err != nil {
		return nil
	}
	return s.s.Reg.Register(addr, srv.ServiceName())
}
