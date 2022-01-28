package endpoints

import "github.com/anqur/gbio/internal/endpoints"

type Option func(s *endpoints.Endpoint)

const DefaultTag = "v1"

func WithTag(tag string) Option {
	return func(s *endpoints.Endpoint) { s.Tag = tag }
}
