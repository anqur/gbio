package gbio

import (
	"bytes"
	"net/http"
	"net/url"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/anqur/gbio/internal/clients"
	"github.com/anqur/gbio/internal/endpoints"
	"github.com/anqur/gbio/internal/registries"
)

var (
	ErrClientEndpointNotGiven = clients.ErrEndpointNotGiven
)

type Client struct {
	Tag string

	cl clients.Client
}

type ClientOption func(c *clients.Client) error

func NewClient(opts ...ClientOption) (c *Client, err error) {
	c = &Client{
		Tag: DefaultTag,
		cl:  clients.Client{H: http.DefaultClient},
	}
	for _, opt := range opts {
		if err = opt(&c.cl); err != nil {
			return
		}
	}
	return
}

func UseClient(opts ...ClientOption) error {
	for _, opt := range opts {
		if err := opt(&DefaultClient.cl); err != nil {
			return err
		}
	}
	return nil
}

var DefaultClient, _ = NewClient(WithEndpoint("http://localhost:8080"))

func WithEndpoint(rawUrl string) ClientOption {
	return func(c *clients.Client) (err error) {
		c.U, err = url.Parse(rawUrl)
		return
	}
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

func WithHttpClient(h *http.Client) ClientOption {
	return func(c *clients.Client) error {
		c.H = h
		return nil
	}
}

func (c *Client) HttpClient() *http.Client { return c.cl.H }

func (c *Client) Request(k, path string, e RequestEncoder) (*http.Request, error) {
	d, ctx, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	u, err := c.cl.LookupEndpoint(k)
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

func (c *Client) Close() error { return c.cl.Close() }

type RequestEncoder interface {
	Marshal() ([]byte, http.Header, error)
}

type ClientEndpoint struct {
	endpoints.Endpoint
	Cl    *Client
	Error error
}
