package loadbalance

import (
	"context"

	"github.com/ibinarytree/koala/registry"
)

type RoundRobinBalance struct {
	index int
}

func (r *RoundRobinBalance) Name() string {
	return "roundrobin"
}

func (r *RoundRobinBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {

	if len(nodes) == 0 {
		err = ErrNotHaveNodes
		return
	}

	r.index = (r.index + 1) % len(nodes)
	node = nodes[r.index]

	return
}
