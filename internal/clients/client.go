package clients

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/anqur/gbio/internal/errors"
	"github.com/anqur/gbio/internal/registries"
)

var (
	ErrEndpointNotGiven = fmt.Errorf("%w: endpoint not given", errors.Err)
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
