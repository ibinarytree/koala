package registry

// 服务抽象
type Service struct {
	Name  string
	Nodes []*Node
}

// 服务节点的抽象
type Node struct {
	Id   string
	IP   string
	Port int
}
