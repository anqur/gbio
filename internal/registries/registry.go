package registries

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/anqur/gbio/internal/errors"
	"github.com/anqur/gbio/internal/loggers"
)

const DefaultPrefix = "gbio-"

var (
	ErrEndpointNotFound = fmt.Errorf("%w: endpoint not found", errors.Err)
)

type Registry struct {
	C      *etcd.Config
	Cp     func() context.Context
	Prefix string

	Cl *etcd.Client
}

func NewRegistry(c *etcd.Config) *Registry {
	return &Registry{
		C:      c,
		Cp:     context.Background,
		Prefix: DefaultPrefix,
	}
}

func (r *Registry) dial() (err error) {
	r.Cl, err = etcd.New(*r.C)
	return
}

// TODO: Put endpoint/service info to etcd.

func (r *Registry) Close() error { return r.Cl.Close() }

type EndpointKey string
type EndpointList []string
type ServiceKey string
type ServiceList []string

type CachedRegistry struct {
	*Registry

	mu sync.RWMutex

	Tick time.Duration
	Lb   LB

	endpoints map[EndpointKey]ServiceList
	services  map[ServiceKey]EndpointList
}

func NewCachedRegistry(c *etcd.Config) *CachedRegistry {
	return &CachedRegistry{
		Registry: NewRegistry(c),

		Tick: 30 * time.Second,
		Lb:   FirstLB(),
	}
}

func (r *CachedRegistry) Close() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.Registry.Close()
}

func (r *CachedRegistry) Started() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.Cl != nil
}

func (r *CachedRegistry) Start() error {
	if err := r.dial(); err != nil {
		return err
	}
	r.fetchOnce()
	go r.runFetch()
	return nil
}

func (r *CachedRegistry) dial() (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.Registry.dial()
}

func (r *CachedRegistry) fetchOnce() {
	r.mu.Lock()
	defer r.mu.Unlock()

	resp, err := r.Cl.Get(r.Cp(), r.Prefix, etcd.WithPrefix())
	if err != nil {
		loggers.Error.Println("Fetch services error:", err)
		return
	}

	r.endpoints = make(map[EndpointKey]ServiceList)
	r.services = make(map[ServiceKey]EndpointList)

	for _, kv := range resp.Kvs {
		ep := strings.TrimPrefix(string(kv.Key), r.Prefix)
		services := strings.Split(string(kv.Value), ",")

		r.endpoints[EndpointKey(ep)] = services
		for _, s := range services {
			k := ServiceKey(s)
			r.services[k] = append(r.services[k], ep)
		}
	}
}

func (r *CachedRegistry) runFetch() {
	for range time.Tick(r.Tick) {
		r.fetchOnce()
	}
}

func (r *CachedRegistry) pick(k ServiceKey) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ep, ok := r.services[k]
	if !ok || len(ep) == 0 {
		return "", fmt.Errorf("%w: %q", ErrEndpointNotFound, k)
	}

	return r.Lb.Pick(ep), nil
}

func (r *CachedRegistry) Lookup(k ServiceKey) (string, error) {
	if !r.Started() {
		if err := r.Start(); err != nil {
			return "", nil
		}
	}
	return r.pick(k)
}
