package lb

import (
	"barton.top/btgo/pkg/discover"
	"sync/atomic"
)

// NewRoundRobin returns a load balancer that returns services in sequence.
func NewRoundRobin(ins *discover.Instancer) Balancer {
	return &roundRobin{
		ins: ins,
		c:   0,
	}
}

type roundRobin struct {
	ins *discover.Instancer
	c   uint64
}

func (rr *roundRobin) Endpoint() (string, error) {
	endpoints := rr.ins.Addresses()
	if len(endpoints) <= 0 {
		return "", ErrNoEndpoints
	}
	old := atomic.AddUint64(&rr.c, 1) - 1
	idx := old % uint64(len(endpoints))
	return endpoints[idx], nil
}
