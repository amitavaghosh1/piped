package loadbalancers

import (
	"math/rand"
	"time"
)

type RandomBalancer struct {
	resourceCount int64
}

func NewRandomdBalancer(resourceCount int) *RandomBalancer {
	return &RandomBalancer{resourceCount: int64(resourceCount)}
}

func (r *RandomBalancer) Get() int64 {
	rand.Seed(time.Now().UnixNano())

	return rand.Int63n(r.resourceCount)
}
