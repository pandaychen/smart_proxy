package wrr

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

//wrr: https://github.com/pandaychen/goes-wrapper/blob/master/balancer/roundrobin.go
//https://github.com/nginx/nginx/commit/52327e0627f49dbda1e8db695e63a4b0af4448b1

// WrrBalancerPool is a instance of backend.BackendNodePool
/*

type BackendNodePool interface {
	//return current pool size
	Size() int

	// get a usable backend nodes from pool
	Pick() *BackendNode

	// add a backend node to pool
	Add(addr string)

	// remove a node from pool
	Remove(addr string)

	// set a node up status
	UpNode(addr string)

	// set a node down status
	DownNode(addr string)
}
*/
type BalancerPool struct {
	sync.RWMutex
	BackendsSet  map[string]struct{}
	BackendLists []*BackendNodeWrapper
	Logger       *zap.Logger
}

func NewBalancerPool(logger *zap.Logger, backends_map map[string]int) (*BalancerPool, error) {
	pool := &BalancerPool{
		Logger:      logger,
		BackendsSet: make(map[string]struct{}),
	}
	for addr, weight := range backends_map {
		pool.Add(addr, weight)
	}
	return pool, nil
}

func (p *BalancerPool) Size() int {
	return len(p.BackendLists)
}

func (p *BalancerPool) Add(addr string, weight int) error {
	defer p.Unlock()
	p.Lock()
	if _, exists := p.BackendsSet[addr]; exists {
		return errors.New("backend exists")
	}

	if weight <= 0 {
		weight = 1
	}

	bnode := CreateWrrBackendNode(addr, weight)
	p.BackendLists == append(p.BackendLists, bnode)
}
