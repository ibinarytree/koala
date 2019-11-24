package middleware

import (
	"context"

	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/meta"
	"github.com/ibinarytree/koala/util"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
)

func TraceServerMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		//从ctx获取grpc的metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			//没有的话,新建一个
			md = metadata.Pairs()
		}

		tracer := opentracing.GlobalTracer()
		parentSpanContext, err := tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(md))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			logs.Warn(ctx, "trace extract failed, parsing trace information: %v", err)
		}

		serverMeta := meta.GetServerMeta(ctx)
		//开始追踪该方法
		serverSpan := tracer.StartSpan(
			serverMeta.Method,
			ext.RPCServerOption(parentSpanContext),
			ext.SpanKindRPCServer,
		)

		serverSpan.SetTag(util.TraceID, logs.GetTraceId(ctx))
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		resp, err = next(ctx, req)
		//记录错误
		if err != nil {
			ext.Error.Set(serverSpan, true)
			serverSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}

		serverSpan.Finish()
		return
	}
}

func TraceRpcMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {

		tracer := opentracing.GlobalTracer()
		var parentSpanCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentSpanCtx = parent.Context()
		}

		opts := []opentracing.StartSpanOption{
			opentracing.ChildOf(parentSpanCtx),
			ext.SpanKindRPCClient,
			opentracing.Tag{Key: string(ext.Component), Value: "koala_rpc"},
			opentracing.Tag{Key: util.TraceID, Value: logs.GetTraceId(ctx)},
		}

		rpcMeta := meta.GetRpcMeta(ctx)
		clientSpan := tracer.StartSpan(rpcMeta.ServiceName, opts...)

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}

		if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, metadataTextMap(md)); err != nil {
			logs.Debug(ctx, "grpc_opentracing: failed serializing trace information: %v", err)
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		ctx = metadata.AppendToOutgoingContext(ctx, util.TraceID, logs.GetTraceId(ctx))
		ctx = opentracing.ContextWithSpan(ctx, clientSpan)

		resp, err = next(ctx, req)
		//记录错误
		if err != nil {
			ext.Error.Set(clientSpan, true)
			clientSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}

		clientSpan.Finish()
		return
	}
}
