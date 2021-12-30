package gbio

import (
	"context"
	"log"
	"time"

	"github.com/anqur/gbio/internal/errors"
	"github.com/anqur/gbio/internal/loggers"
	"github.com/anqur/gbio/internal/registries"
)

var (
	Err = errors.Err

	ErrRegistryEndpointNotFound = registries.ErrEndpointNotFound
	ErrRegistryEmptyServiceInfo = registries.ErrEmptyServiceInfo
)

func UseErrorLogger(l *log.Logger) { loggers.Error = l }

type (
	LookupRegistryOption func(r *registries.CachedRegistry)
	ContextProvider      = func() context.Context
)

func WithTick(d time.Duration) LookupRegistryOption {
	return func(r *registries.CachedRegistry) { r.Tick = d }
}

func WithLookupContext(f ContextProvider) LookupRegistryOption {
	return func(r *registries.CachedRegistry) { r.Cp = f }
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

func WithServiceContext(f ContextProvider) ServiceRegistryOption {
	return func(r *registries.Registry) { r.Cp = f }
}

func WithServicePrefix(p string) ServiceRegistryOption {
	return func(r *registries.Registry) { r.Prefix = p }
}
