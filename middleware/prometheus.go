package middleware

import (
	"context"
	"time"
	"net/http"
	"fmt"

	"github.com/ibinarytree/koala/meta"
	"github.com/ibinarytree/koala/middleware/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	DefaultServerMetrics = prometheus.NewServerMetrics()
	DefaultRpcMetrics    = prometheus.NewRpcMetrics()
)

func init() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		addr := fmt.Sprintf("0.0.0.0:%d", 8888)
		http.ListenAndServe(addr, nil)
	}()
}

func PrometheusServerMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {

		serverMeta := meta.GetServerMeta(ctx)
		DefaultServerMetrics.IncrRequest(ctx, serverMeta.ServiceName, serverMeta.Method)

		startTime := time.Now()
		resp, err = next(ctx, req)

		DefaultServerMetrics.IncrCode(ctx, serverMeta.ServiceName, serverMeta.Method, err)
		DefaultServerMetrics.Latency(ctx, serverMeta.ServiceName,
			serverMeta.Method, time.Since(startTime).Nanoseconds()/1000)
		return
	}
}

func PrometheusRpcMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {

		rpcMeta := meta.GetRpcMeta(ctx)
		DefaultRpcMetrics.IncrRequest(ctx, rpcMeta.ServiceName, rpcMeta.Method)

		startTime := time.Now()
		resp, err = next(ctx, req)

		DefaultRpcMetrics.IncrCode(ctx, rpcMeta.ServiceName, rpcMeta.Method, err)
		DefaultRpcMetrics.Latency(ctx, rpcMeta.ServiceName,
			rpcMeta.Method, time.Since(startTime).Nanoseconds()/1000)
		return
	}
}
