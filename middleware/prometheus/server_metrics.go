package prometheus

import (
	"context"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/status"
)

// koala服务端采样打点
type ServerMetrics struct {
	serverRequestCounter *prom.CounterVec
	serverCodeCounter    *prom.CounterVec
	serverLatencySummary *prom.SummaryVec
}

//生成server metrics实例
func NewServerMetrics() *ServerMetrics {
	return &ServerMetrics{
		serverRequestCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "koala_server_request_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, []string{"service", "method"}),
		serverCodeCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "koala_server_handled_code_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, []string{"service", "method", "grpc_code"}),
		serverLatencySummary: promauto.NewSummaryVec(
			prom.SummaryOpts{
				Name:       "koala_proc_cost",
				Help:       "RPC latency distributions.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service", "method"},
		),
	}
}

func (m *ServerMetrics) IncrRequest(ctx context.Context, serviceName, methodName string) {
	m.serverRequestCounter.WithLabelValues(serviceName, methodName).Inc()
}

func (m *ServerMetrics) IncrCode(ctx context.Context, serviceName, methodName string, err error) {
	st, _ := status.FromError(err)
	m.serverCodeCounter.WithLabelValues(serviceName, methodName, st.Code().String()).Inc()
}

func (m *ServerMetrics) Latency(ctx context.Context, serviceName, methodName string, us int64) {

	m.serverLatencySummary.WithLabelValues(serviceName, methodName).Observe(float64(us))
}
