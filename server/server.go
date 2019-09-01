package server

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ibinarytree/koala/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

type KoalaServer struct {
	*grpc.Server
	limiter *rate.Limiter

	userMiddleware []middleware.Middleware
}

var koalaServer = &KoalaServer{
	Server: grpc.NewServer(),
}

func Use(m ...middleware.Middleware) {
	koalaServer.userMiddleware = append(koalaServer.userMiddleware, m...)
}

func Init(serviceName string) (err error) {
	err = InitConfig(serviceName)
	if err != nil {
		return
	}

	//初始化限流器
	if koalaConf.Limit.SwitchOn {
		koalaServer.limiter = rate.NewLimiter(rate.Limit(koalaConf.Limit.QPSLimit),
			koalaConf.Limit.QPSLimit)
	}
	return
}

func Run() {
	if koalaConf.Prometheus.SwitchOn {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			addr := fmt.Sprintf("0.0.0.0:%d", koalaConf.Prometheus.Port)
			log.Fatal(http.ListenAndServe(addr, nil))
		}()
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", koalaConf.Port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	koalaServer.Serve(lis)
}

func GRPCServer() *grpc.Server {
	return koalaServer.Server
}

func BuildServerMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware
	if koalaConf.Prometheus.SwitchOn {
		mids = append(mids, middleware.PrometheusServerMiddleware)
	}

	if koalaConf.Limit.SwitchOn {
		mids = append(mids, middleware.NewRateLimitMiddleware(koalaServer.limiter))
	}

	if len(koalaServer.userMiddleware) != 0 {
		mids = append(mids, koalaServer.userMiddleware...)
	}

	if len(mids) > 0 {
		//把所有中间件组织成一个调用链
		m := middleware.Chain(mids[0], mids[1:]...)
		//返回调用链的入口函数
		return m(handle)
	}

	return handle
}
