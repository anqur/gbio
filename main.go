package gbio

import (
	"github.com/anqur/gbio/internal/clients"
)

type Client struct {
	s clients.Setting
}

type Option interface {
	Apply(s *clients.Setting)
}

func New(opts ...Option) (c *Client) {
	c = new(Client)
	for _, opt := range opts {
		opt.Apply(&c.s)
	}
	return
}
