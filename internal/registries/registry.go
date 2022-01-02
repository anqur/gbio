package registries

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/anqur/gbio/internal/errors"
	"github.com/anqur/gbio/internal/loggers"
)

const DefaultPrefix = "gbio-"

var (
	ErrEndpointNotFound = fmt.Errorf("%w: endpoint not found", errors.Err)
	ErrEmptyServiceInfo = fmt.Errorf("%w: empty service info", errors.Err)
)

type Registry struct {
	C      *etcd.Config
	Cp     func() context.Context
	Prefix string
	Tick   time.Duration

	cl *etcd.Client
	tx *concurrency.Session
}

func NewRegistry(c *etcd.Config) *Registry {
	return &Registry{
		C:      c,
		Cp:     context.Background,
		Prefix: DefaultPrefix,
		Tick:   30 * time.Second,
	}
}

func (r *Registry) dial() (err error) {
	r.cl, err = etcd.New(*r.C)
	if err != nil {
		return
	}
	ttl := int(r.Tick.Seconds())
	r.tx, err = concurrency.NewSession(r.cl, concurrency.WithTTL(ttl))
	return
}

func (r *Registry) Register(serviceKey, serviceName string) error {
	if serviceKey == "" || serviceName == "" {
		return fmt.Errorf(
			"%w: key=%q, name=%q",
			ErrEmptyServiceInfo,
			serviceKey,
			serviceName,
		)
	}
	prefixedKey := r.Prefix + serviceKey
	_, err := r.tx.
		Client().
		Put(
			r.Cp(),
			prefixedKey,
			serviceName,
			etcd.WithLease(r.tx.Lease()),
		)
	return err
}

func (r *Registry) Close() error {
	_ = r.tx.Close()
	return r.cl.Close()
}

type EndpointKey string
type EndpointList []string
type ServiceKey string

type CachedRegistry struct {
	*Registry

	mu        sync.RWMutex
	Lb        LB
	endpoints map[EndpointKey]ServiceKey
	services  map[ServiceKey]EndpointList
}

func NewCachedRegistry(c *etcd.Config) *CachedRegistry {
	return &CachedRegistry{
		Registry: NewRegistry(c),
		Lb:       FirstLB(),
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

	return r.cl != nil
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

	resp, err := r.cl.Get(r.Cp(), r.Prefix, etcd.WithPrefix())
	if err != nil {
		loggers.Error.Println("Fetch services error:", err)
		return
	}

	r.endpoints = make(map[EndpointKey]ServiceKey)
	r.services = make(map[ServiceKey]EndpointList)

	for _, kv := range resp.Kvs {
		ep := strings.TrimPrefix(string(kv.Key), r.Prefix)
		srv := ServiceKey(kv.Value)
		r.endpoints[EndpointKey(ep)] = srv
		r.services[srv] = append(r.services[srv], ep)
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
