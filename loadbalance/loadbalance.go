package loadbalance

import (
	"context"
	"errors"

	"github.com/ibinarytree/koala/registry"
)

var (
	ErrNotHaveNodes = errors.New("not have node")
)

const (
	DefaultNodeWeight = 100
)

type LoadBalanceType int

const (
	LoadBalanceTypeRandom = iota
	LoadBalanceTypeRoundRobin 
)

type LoadBalance interface {
	Name() string
	Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error)
}

func GetLoadBalance(balanceType LoadBalanceType) (balancer LoadBalance) {

	switch (balanceType) {
		 case LoadBalanceTypeRandom:
			 balancer = NewRandomBalance()
		 case LoadBalanceTypeRoundRobin:
			 balancer = NewRoundRobinBalance()
		 default:
			 balancer = NewRandomBalance()
	}
	return
}
