package loadbalancers

import (
	"sync/atomic"
)

type Balancer interface {
	Get() int64
}

type RoundRobalancer struct {
	resourceCount int64
	lastIndex     *int64
}

func NewRoundRobalancer(resourceCount int) *RoundRobalancer {
	var index = int64(0)
	return &RoundRobalancer{resourceCount: int64(resourceCount), lastIndex: &index}
}

func (r *RoundRobalancer) Get() int64 {
	idx := atomic.LoadInt64(r.lastIndex)

	if idx < r.resourceCount-1 {
		atomic.AddInt64(r.lastIndex, 1)
	} else {
		atomic.StoreInt64(r.lastIndex, 0)
	}

	return idx
}
