package discov

import (
	"google.golang.org/grpc/resolver"
	"sync"
)

type Resolver struct {
	sync.RWMutex
	cc resolver.ClientConn
}

func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *Resolver) Close() {}

func (r *Resolver) Update(addrs []string) {
	addresses := r.getAddresses(addrs)
	r.cc.UpdateState(resolver.State{Addresses: addresses})
}

func (r *Resolver) getAddresses(addrs []string) []resolver.Address {
	addresses := make([]resolver.Address, len(addrs))
	for _, addr := range addrs {
		addresses = append(addresses, resolver.Address{
			Addr: addr,
		})
	}
	return addresses
}
