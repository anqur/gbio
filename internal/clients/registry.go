package clients

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

var (
	ErrEndpointNotFound = fmt.Errorf("%w: endpoint not found", errors.Err)
)

type EndpointKey string
type EndpointList []string
type ServiceKey string
type ServiceList []string

type Registry struct {
	mu sync.RWMutex

	C      *etcd.Config
	Tick   time.Duration
	Cp     func() context.Context
	Prefix string
	Lb     LB

	cl        *etcd.Client
	endpoints map[EndpointKey]ServiceList
	services  map[ServiceKey]EndpointList
}

func NewRegistry(c *etcd.Config) *Registry {
	return &Registry{
		C:      c,
		Tick:   30 * time.Second,
		Cp:     context.Background,
		Prefix: "gbio-",
		Lb:     FirstLB(),
	}
}

func (r *Registry) Close() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.cl.Close()
}

func (r *Registry) Started() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.cl != nil
}

func (r *Registry) Start() error {
	if err := r.dial(); err != nil {
		return err
	}
	r.fetchOnce()
	go r.runFetch()
	return nil
}

func (r *Registry) dial() (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cl, err = etcd.New(*r.C)
	return
}

func (r *Registry) fetchOnce() {
	r.mu.Lock()
	defer r.mu.Unlock()

	resp, err := r.cl.Get(r.Cp(), r.Prefix, etcd.WithPrefix())
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

func (r *Registry) runFetch() {
	for range time.Tick(r.Tick) {
		r.fetchOnce()
	}
}

func (r *Registry) pick(k ServiceKey) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ep, ok := r.services[k]
	if !ok || len(ep) == 0 {
		return "", fmt.Errorf("%w: %q", ErrEndpointNotFound, k)
	}

	return r.Lb.Pick(ep), nil
}

func (r *Registry) PickUpstream(k ServiceKey) (string, error) {
	if !r.Started() {
		if err := r.Start(); err != nil {
			return "", nil
		}
	}
	return r.pick(k)
}
