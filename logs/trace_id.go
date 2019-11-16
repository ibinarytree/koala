package logs

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const (
	MaxTraceId = 100000000
)

type traceIdKey struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetTraceId(ctx context.Context) (traceId string) {

	traceId, ok := ctx.Value(traceIdKey{}).(string)
	if !ok {
		traceId = GenTraceId()
	}

	return
}

func GenTraceId() (traceId string) {
	now := time.Now()
	traceId = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%08d", now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), rand.Int31n(MaxTraceId))
	return
}

func WithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, traceIdKey{}, traceId)
}
