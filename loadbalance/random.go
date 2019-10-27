package loadbalance

import (
	"context"
	"math/rand"

	"github.com/ibinarytree/koala/errno"
	"github.com/ibinarytree/koala/registry"
)

type RandomBalance struct {
	name string
}

func NewRandomBalance() LoadBalance {
	return &RandomBalance{
		name: "random",
	}
}

func (r *RandomBalance) Name() string {
	return r.name
}

func (r *RandomBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {

	if len(nodes) == 0 {
		err = errno.NotHaveInstance
		return
	}

	defer func() {
		if node != nil {
			setSelected(ctx, node)
		}
	}()

	var newNodes  = filterNodes(ctx, nodes)
	if len(newNodes) == 0 {
		err = errno.AllNodeFailed
		return
	}

	var totalWeight int
	for _, val := range newNodes {
		if val.Weight == 0 {
			val.Weight = DefaultNodeWeight
		}
		totalWeight += val.Weight
	}

	curWeight := rand.Intn(totalWeight)
	curIndex := -1
	for index, node := range nodes {
		curWeight -= node.Weight
		if curWeight < 0 {
			curIndex = index
			break
		}
	}

	if curIndex == -1 {
		err = errno.AllNodeFailed
		return
	}

	node = nodes[curIndex]
	return
}
