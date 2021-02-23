package robin

import (
	"github.com/mhchlib/go-kit/endpoint"
	"github.com/mhchlib/go-kit/sd"
	"github.com/mhchlib/go-kit/sd/lb"
	"sync/atomic"
)

func NewListRobin(s sd.Endpointer) lb.Balancer {
	return &ListRobin{
		s: s,
		c: 0,
	}
}

type ListRobin struct {
	s sd.Endpointer
	c uint64
}

func (rr *ListRobin) Endpoint() (endpoint.Endpoint, error) {
	endpoints, err := rr.s.Endpoints()
	if err != nil {
		return nil, err
	}
	if len(endpoints) <= 0 {
		return nil, lb.ErrNoEndpoints
	}
	old := atomic.AddUint64(&rr.c, 1) - 1
	idx := old % uint64(len(endpoints))
	return endpoints[idx], nil
}

func (rr *ListRobin) Endpoints() ([]endpoint.Endpoint, error) {
	endpoints, err := rr.s.Endpoints()
	if err != nil {
		return nil, err
	}
	if len(endpoints) <= 0 {
		return nil, lb.ErrNoEndpoints
	}
	return endpoints, nil
}
