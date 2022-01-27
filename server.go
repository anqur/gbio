package gbio

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/anqur/gbio/logging"
	etcd "go.etcd.io/etcd/client/v3"

	"github.com/anqur/gbio/internal/endpoints"
	"github.com/anqur/gbio/internal/registries"
	"github.com/anqur/gbio/internal/utils"
)

type Server struct {
	u   *url.URL
	m   *http.ServeMux
	s   http.Server
	eps map[string]*ServerEndpoint
	reg *registries.Registry
}

type ServerOption func(s *Server)

func NewServer(rawURL string, opts ...ServerOption) (*Server, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	m := new(http.ServeMux)
	s := &Server{
		u: u,
		m: m,
		s: http.Server{
			Addr:    u.Host,
			Handler: m,
		},
		eps: make(map[string]*ServerEndpoint),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func UseServer(opts ...ServerOption) {
	for _, opt := range opts {
		opt(DefaultServer)
	}
}

var DefaultServer, _ = NewServer("http://0.0.0.0:8080")

func WithServiceRegistry(
	cfg *etcd.Config,
	opts ...ServiceRegistryOption,
) ServerOption {
	return func(s *Server) {
		s.reg = registries.NewRegistry(cfg)
		for _, opt := range opts {
			opt(s.reg)
		}
	}
}

func ListenAndServe() error { return DefaultServer.ListenAndServe() }

func (s *Server) ListenAndServe() error {
	for _, srv := range s.eps {
		logging.Info.Println("Registering:", srv.Name, srv.BaseURI)
		s.m.HandleFunc(srv.BaseURI, srv.Handler)
	}
	if s.reg != nil {
		if err := s.registerServer(); err != nil {
			return err
		}
	}
	logging.Info.Println("Listening:", s.s.Addr)
	return s.s.ListenAndServe()
}

func (s *Server) Register(srv *ServerEndpoint) {
	s.eps[srv.Name] = srv
}

func (s *Server) serverAddr() (string, error) {
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

func (s *Server) registerServer() error {
	addr, err := s.serverAddr()
	if err != nil {
		return nil
	}
	var eps []*endpoints.Endpoint
	for _, ep := range s.eps {
		eps = append(eps, &ep.Endpoint)
	}
	return s.reg.Register(addr, eps)
}

type ServerEndpoint struct {
	endpoints.Endpoint

	Handler http.HandlerFunc
}
