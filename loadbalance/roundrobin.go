package loadbalance

import (
	"context"

	"github.com/ibinarytree/koala/errno"
	"github.com/ibinarytree/koala/registry"
)

type RoundRobinBalance struct {
	name  string
	index int
}

func NewRoundRobinBalance() LoadBalance {
	return &RoundRobinBalance{
		name: "roundrobin",
	}
}

func (r *RoundRobinBalance) Name() string {
	return r.name
}

func (r *RoundRobinBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {

	if len(nodes) == 0 {
		err = errno.NotHaveInstance
		return
	}

	defer func() {
		if node != nil {
			setSelected(ctx, node)
		}
	}()

	var newNodes = filterNodes(ctx, nodes)
	if len(newNodes) == 0 {
		err = errno.AllNodeFailed
		return
	}

	r.index = (r.index + 1) % len(nodes)
	node = nodes[r.index]
	return
}
