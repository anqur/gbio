package gbio

import (
	"net/http"
	"net/url"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/anqur/gbio/internal/clients"
	"github.com/anqur/gbio/internal/registries"
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

func WithLookupRegistry(
	cfg *etcd.Config,
	opts ...LookupRegistryOption,
) ClientOption {
	return func(c *clients.Client) error {
		c.Reg = registries.NewCachedRegistry(cfg)
		for _, opt := range opts {
			opt(c.Reg)
		}
		return nil
	}
}
func UseLookupRegistry(c *etcd.Config, opts ...LookupRegistryOption) {
	DefaultClient.cl.Reg = registries.NewCachedRegistry(c)
	for _, opt := range opts {
		opt(DefaultClient.cl.Reg)
	}
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
