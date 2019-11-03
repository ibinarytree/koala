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
	options  []*proto.Option

	//服务唯一标识，用来服务注册以及发现，用点进行分隔，比如 google.gmail.account.user
	ServiceName string
	//serviceName以点分隔而成的数组
	serviceNameParts []string
	//serviceName转成路径的格式,比如:google.gmail.account.user转成google/gmail/account/user
	ServiceNamePartsPath string
	//包名
	PackageName string
	//是否包含go_package
	containGoPackage bool
	Prefix           string
}
