package gbio

import (
	"context"
	"net/http"
	"net/url"
	"time"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/anqur/gbio/internal/clients"
)

var (
	ErrClientEndpointNotGiven = clients.ErrEndpointNotGiven
)

type Client struct {
	cl clients.Client
}

type ClientOption func(c *clients.Client) error

func NewClient(opts ...ClientOption) (c *Client, err error) {
	c = new(Client)
	for _, opt := range opts {
		if err = opt(&c.cl); err != nil {
			return
		}
	}
	return
}

var DefaultClient, _ = NewClient(
	WithHttpClient(http.DefaultClient),
	WithEndpoint("http://localhost:8080"),
)

func WithEndpoint(rawUrl string) ClientOption {
	return func(c *clients.Client) (err error) {
		c.U, err = url.Parse(rawUrl)
		return
	}
}
func UseEndpoint(rawURL string) (err error) {
	DefaultClient.cl.U, err = url.Parse(rawURL)
	return
}

func WithRegistry(cfg *etcd.Config, opts ...RegistryOption) ClientOption {
	return func(c *clients.Client) error {
		c.Reg = clients.NewRegistry(cfg)
		for _, opt := range opts {
			opt(c.Reg)
		}
		return nil
	}
}
func UseRegistry(c *etcd.Config, opts ...RegistryOption) {
	DefaultClient.cl.Reg = clients.NewRegistry(c)
	for _, opt := range opts {
		opt(DefaultClient.cl.Reg)
	}
}

type RegistryOption func(r *clients.Registry)

func WithTick(d time.Duration) RegistryOption {
	return func(r *clients.Registry) { r.Tick = d }
}

type ContextProvider = func() context.Context

func WithContext(f ContextProvider) RegistryOption {
	return func(r *clients.Registry) { r.Cp = f }
}

func WithPrefix(p string) RegistryOption {
	return func(r *clients.Registry) { r.Prefix = p }
}

func WithPickFirst() RegistryOption {
	return func(r *clients.Registry) { r.Lb = clients.FirstLB() }
}

func WithPickRandom() RegistryOption {
	return func(r *clients.Registry) { r.Lb = clients.RandLB() }
}

func WithHttpClient(h *http.Client) ClientOption {
	return func(c *clients.Client) error {
		c.H = h
		return nil
	}
}
func UseHttpClient(h *http.Client) { DefaultClient.cl.H = h }

func (c *Client) HttpClient() *http.Client { return c.cl.H }

func (c *Client) LookupEndpoint(serviceKey string) (string, error) {
	return c.cl.LookupEndpoint(serviceKey)
}

func (c *Client) Close() error { return c.cl.Close() }
