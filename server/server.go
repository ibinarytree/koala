package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/middleware"
	"github.com/ibinarytree/koala/registry"
	_ "github.com/ibinarytree/koala/registry/etcd"
	"github.com/ibinarytree/koala/util"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

type KoalaServer struct {
	*grpc.Server
	limiter        *rate.Limiter
	register       registry.Registry
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

	initLogger()

	//初始化注册中心
	err = initRegister(serviceName)
	if err != nil {
		logs.Error(context.TODO(), "init register failed, err:%v", err)
		return
	}

	err = initTrace(serviceName)
	if err != nil {
		logs.Error(context.TODO(), "init tracing failed, err:%v", err)
	}
	return
}

func initTrace(serviceName string) (err error) {

	if !koalaConf.Trace.SwitchOn {
		return
	}

	return middleware.InitTrace(serviceName, koalaConf.Trace.ReportAddr,
		koalaConf.Trace.SampleType, koalaConf.Trace.SampleRate)
}

func initLogger() (err error) {
	filename := fmt.Sprintf("%s/%s.log", koalaConf.Log.Dir, koalaConf.ServiceName)
	outputer, err := logs.NewFileOutputer(filename)
	if err != nil {
		return
	}

	level := logs.GetLogLevel(koalaConf.Log.Level)
	logs.InitLogger(level, koalaConf.Log.ChanSize, koalaConf.ServiceName)
	logs.AddOutputer(outputer)

	if koalaConf.Log.ConsoleLog {
		logs.AddOutputer(logs.NewConsoleOutputer())
	}
	return
}

func initRegister(serviceName string) (err error) {

	if !koalaConf.Regiser.SwitchOn {
		return
	}

	ctx := context.TODO()
	registryInst, err := registry.InitRegistry(ctx,
		koalaConf.Regiser.RegisterName,
		registry.WithAddrs([]string{koalaConf.Regiser.RegisterAddr}),
		registry.WithTimeout(koalaConf.Regiser.Timeout),
		registry.WithRegistryPath(koalaConf.Regiser.RegisterPath),
		registry.WithHeartBeat(koalaConf.Regiser.HeartBeat),
	)
	if err != nil {
		logs.Error(ctx, "init registry failed, err:%v", err)
		return
	}

	koalaServer.register = registryInst
	service := &registry.Service{
		Name: serviceName,
	}

	ip, err := util.GetLocalIP()
	if err != nil {
		return
	}
	service.Nodes = append(service.Nodes, &registry.Node{
		IP:   ip,
		Port: koalaConf.Port,
	},
	)

	registryInst.Register(context.TODO(), service)
	return
}

func Run() {
	/*
		if koalaConf.Prometheus.SwitchOn {
			go func() {
				http.Handle("/metrics", promhttp.Handler())
				addr := fmt.Sprintf("0.0.0.0:%d", koalaConf.Prometheus.Port)
				log.Fatal(http.ListenAndServe(addr, nil))
			}()
		}*/

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

	mids = append(mids, middleware.AccessLogMiddleware)
	if koalaConf.Prometheus.SwitchOn {
		mids = append(mids, middleware.PrometheusServerMiddleware)
	}

	if koalaConf.Limit.SwitchOn {
		mids = append(mids, middleware.NewRateLimitMiddleware(koalaServer.limiter))
	}

	if koalaConf.Trace.SwitchOn {
		mids = append(mids, middleware.TraceServerMiddleware)
	}

	if len(koalaServer.userMiddleware) != 0 {
		mids = append(mids, koalaServer.userMiddleware...)
	}

	m := middleware.Chain(middleware.PrepareMiddleware, mids...)
	return m(handle)
}
