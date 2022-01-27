package gbio

import (
	"context"
	"fmt"
	"time"

	"github.com/anqur/gbio/internal/endpoints"
	"github.com/anqur/gbio/internal/errors"
	"github.com/anqur/gbio/internal/registries"
)

var (
	Err = errors.Err

	ErrCodecBadMsgTag  = fmt.Errorf("%w: unknown message tag", errors.Err)
	ErrCodecBadMsgType = fmt.Errorf("%w: unknown message type", errors.Err)

	ErrRegistryEndpointNotFound = registries.ErrEndpointNotFound
	ErrRegistryEmptyServiceInfo = registries.ErrEmptyServiceInfo
)

type LookupRegistryOption func(r *registries.CachedRegistry)

func WithTick(d time.Duration) LookupRegistryOption {
	return func(r *registries.CachedRegistry) { r.Tick = d }
}

func WithLookupContext(ctx context.Context) LookupRegistryOption {
	return func(r *registries.CachedRegistry) { r.Ctx = ctx }
}

func WithLookupPrefix(p string) LookupRegistryOption {
	return func(r *registries.CachedRegistry) { r.Prefix = p }
}

func WithPickFirst() LookupRegistryOption {
	return func(r *registries.CachedRegistry) { r.Lb = registries.FirstLB() }
}

func WithPickRandom() LookupRegistryOption {
	return func(r *registries.CachedRegistry) { r.Lb = registries.RandLB() }
}

type ServiceRegistryOption func(r *registries.Registry)

func WithServiceContext(ctx context.Context) ServiceRegistryOption {
	return func(r *registries.Registry) { r.Ctx = ctx }
}

func WithServicePrefix(p string) ServiceRegistryOption {
	return func(r *registries.Registry) { r.Prefix = p }
}

type EndpointOption func(s *endpoints.Endpoint)

const DefaultTag = "v1"

func WithTag(tag string) EndpointOption {
	return func(s *endpoints.Endpoint) { s.Tag = tag }
}
