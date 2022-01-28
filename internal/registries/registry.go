package registries

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/anqur/gbio/logging"
	"github.com/anqur/gbio/registries"
	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/anqur/gbio/internal/endpoints"
)

const DefaultPrefix = "gbio-"

type Registry struct {
	C      *etcd.Config
	Ctx    context.Context
	Prefix string
	Tick   time.Duration

	cl *etcd.Client
	tx *concurrency.Session
}

func New(c *etcd.Config) *Registry {
	return &Registry{
		C:      c,
		Ctx:    context.Background(),
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

func (r *Registry) Register(addr string, eps []*endpoints.Endpoint) error {
	if addr == "" || eps == nil {
		return fmt.Errorf(
			"%w: addr=%q, endpoints=%q",
			registries.ErrEmptyEndpoints,
			addr,
			eps,
		)
	}
	vals, err := json.Marshal(eps)
	if err != nil {
		return err
	}
	_, err = r.tx.
		Client().
		Put(
			r.Ctx,
			r.Prefix+addr,
			string(vals),
			etcd.WithLease(r.tx.Lease()),
		)
	return err
}

func (r *Registry) Close() error {
	_ = r.tx.Close()
	return r.cl.Close()
}

type (
	NodeAddr     string
	NodeList     []string
	EndpointName string
)

type CachedRegistry struct {
	*Registry

	mu    sync.RWMutex
	Lb    LB
	nodes map[NodeAddr][]*endpoints.Endpoint
	eps   map[EndpointName]NodeList
}

func NewCached(c *etcd.Config) *CachedRegistry {
	return &CachedRegistry{
		Registry: New(c),
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

	resp, err := r.cl.Get(r.Ctx, r.Prefix, etcd.WithPrefix())
	if err != nil {
		logging.Error.Println("Fetch endpoints error:", err)
		return
	}

	r.nodes = make(map[NodeAddr][]*endpoints.Endpoint)
	r.eps = make(map[EndpointName]NodeList)

	for _, kv := range resp.Kvs {
		node := strings.TrimPrefix(string(kv.Key), r.Prefix)
		var eps []*endpoints.Endpoint
		if err := json.Unmarshal(kv.Value, &eps); err != nil {
			logging.Error.Println("Unmarshal endpoints error:", err)
			return
		}
		r.nodes[NodeAddr(node)] = eps
		for _, ep := range eps {
			n := EndpointName(ep.Name)
			r.eps[n] = append(r.eps[n], node)
		}
	}
}

func (r *CachedRegistry) runFetch() {
	for range time.Tick(r.Tick) {
		r.fetchOnce()
	}
}

func (r *CachedRegistry) pick(k EndpointName) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	eps, ok := r.eps[k]
	if !ok || len(eps) == 0 {
		return "", fmt.Errorf("%w: %q", registries.ErrEndpointNotFound, k)
	}
	if len(eps) == 1 {
		return eps[0], nil
	}

	return r.Lb.Pick(eps), nil
}

func (r *CachedRegistry) Lookup(k EndpointName) (string, error) {
	if !r.Started() {
		if err := r.Start(); err != nil {
			return "", nil
		}
	}
	return r.pick(k)
}
