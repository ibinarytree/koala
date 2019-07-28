package main

import (
	"github.com/ibinarytree/proto"
)

type Generator interface {
	Run(opt *Option, metaData *ServiceMetaData) error
}

type ServiceMetaData struct {
	Service  *proto.Service
	Messages []*proto.Message
	Rpc      []*proto.RPC
	Package  *proto.Package
	Prefix   string
}
