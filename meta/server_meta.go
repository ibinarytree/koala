package meta

import (
	"context"
)

type ServerMeta struct {
	ServiceName string
	Method      string
	Cluster     string
	TraceID     string
	Env         string
	ServerIP    string
	ClientIP    string
	IDC         string
}

type serverMetaContextKey struct{}

func GetServerMeta(ctx context.Context) *ServerMeta {
	meta, ok := ctx.Value(serverMetaContextKey{}).(*ServerMeta)
	if !ok {
		meta = &ServerMeta{}
	}

	return meta
}

func InitServerMeta(ctx context.Context, service, method string) context.Context {
	meta := &ServerMeta{
		Method:      method,
		ServiceName: service,
	}
	return context.WithValue(ctx, serverMetaContextKey{}, meta)
}
