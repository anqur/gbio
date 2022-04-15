package clients

import (
	"bytes"
	"context"
	"fmt"
	"github.com/anqur/gbio/pkg/gbioerr"
	"github.com/anqur/gbio/pkg/encoding"
	"net/http"
	"net/url"
	"time"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/anqur/gbio/internal/endpoints"
	"github.com/anqur/gbio/internal/registries"
)

var (
	ErrEndpointNotGiven = fmt.Errorf("%w: endpoint not given", gbioerr.Err)
)

type Client struct {
	h   *http.Client
	u   *url.URL
	reg *registries.CachedRegistry
}

func (c *Client) lookup(epName string) (string, error) {
	if c.u != nil {
		return (&url.URL{Scheme: c.u.Scheme, Host: c.u.Host}).String(), nil
	}
	if c.reg != nil {
		return c.reg.Lookup(registries.EndpointName(epName))
	}
	return "", ErrEndpointNotGiven
}

func (c *Client) Close() error {
	if c.reg == nil {
		return nil
	}
	return c.reg.Close()
}

type Option func(c *Client) error

func New(opts ...Option) (c *Client, err error) {
	c = &Client{h: http.DefaultClient}
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
		c.u, err = url.Parse(rawUrl)
		return
	}
}

func WithRegistry(
	cfg *etcd.Config,
	opts ...RegistryOption,
) Option {
	return func(c *Client) error {
		c.reg = registries.NewCached(cfg)
		for _, opt := range opts {
			opt(c.reg)
		}
		return nil
	}
}

func WithHttpClient(h *http.Client) Option {
	return func(c *Client) error {
		c.h = h
		return nil
	}
}

func (c *Client) HttpClient() *http.Client { return c.h }

func (c *Client) Request(
	epName, path string,
	e encoding.Encoder,
) (*http.Request, error) {
	d, ctx, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	u, err := c.lookup(epName)
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
