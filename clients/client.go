package clients

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/anqur/gbio/base"
	"github.com/anqur/gbio/codec"
	etcd "go.etcd.io/etcd/client/v3"

	"github.com/anqur/gbio/internal/endpoints"
	"github.com/anqur/gbio/internal/registries"
)

var (
	ErrEndpointNotGiven = fmt.Errorf("%w: endpoint not given", base.Err)
)

type Client struct {
	H   *http.Client
	U   *url.URL
	Reg *registries.CachedRegistry
}

func (c *Client) LookupEndpoint(serviceKey string) (string, error) {
	if c.U != nil {
		return (&url.URL{Scheme: c.U.Scheme, Host: c.U.Host}).String(), nil
	}
	if c.Reg != nil {
		return c.Reg.Lookup(registries.EndpointName(serviceKey))
	}
	return "", ErrEndpointNotGiven
}

func (c *Client) Close() error {
	return c.Reg.Close()
}

type Option func(c *Client) error

func New(opts ...Option) (c *Client, err error) {
	c = &Client{H: http.DefaultClient}
	for _, opt := range opts {
		if err = opt(c); err != nil {
			return
		}
	}
	return
}

func Use(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(Default); err != nil {
			return err
		}
	}
	return nil
}

var Default, _ = New(WithURL("http://localhost:8080"))

func WithURL(rawUrl string) Option {
	return func(c *Client) (err error) {
		c.U, err = url.Parse(rawUrl)
		return
	}
}

func WithRegistry(
	cfg *etcd.Config,
	opts ...RegistryOption,
) Option {
	return func(c *Client) error {
		c.Reg = registries.NewCached(cfg)
		for _, opt := range opts {
			opt(c.Reg)
		}
		return nil
	}
}

func WithHttpClient(h *http.Client) Option {
	return func(c *Client) error {
		c.H = h
		return nil
	}
}

func (c *Client) HttpClient() *http.Client { return c.H }

func (c *Client) Request(k, path string, e codec.Encoder) (*http.Request, error) {
	d, ctx, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	u, err := c.LookupEndpoint(k)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u+path, bytes.NewReader(d))
	if err != nil {
		return nil, err
	}
	for k, vs := range ctx {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	return req, nil
}

type Endpoint struct {
	endpoints.Endpoint
	Cl    *Client
	Error error
}

type RegistryOption func(r *registries.CachedRegistry)

func WithTick(d time.Duration) RegistryOption {
	return func(r *registries.CachedRegistry) { r.Tick = d }
}

func WithContext(ctx context.Context) RegistryOption {
	return func(r *registries.CachedRegistry) { r.Ctx = ctx }
}

func WithPrefix(p string) RegistryOption {
	return func(r *registries.CachedRegistry) { r.Prefix = p }
}

func WithPickFirst() RegistryOption {
	return func(r *registries.CachedRegistry) { r.Lb = registries.FirstLB() }
}

func WithPickRandom() RegistryOption {
	return func(r *registries.CachedRegistry) { r.Lb = registries.RandLB() }
}
