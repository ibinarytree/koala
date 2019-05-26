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

type LoadBalance interface {
	Name() string
	Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error)
}
