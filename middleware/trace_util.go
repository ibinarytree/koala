package middleware

import (
	"context"
	"encoding/base64"
	"strings"

	"fmt"

	"github.com/ibinarytree/koala/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"google.golang.org/grpc/metadata"
)

const (
	binHdrSuffix = "-bin"
)

// metadataTextMap extends a metadata.MD to be an opentracing textmap
type metadataTextMap metadata.MD

// Set is a opentracing.TextMapReader interface that extracts values.
func (m metadataTextMap) Set(key, val string) {
	// gRPC allows for complex binary values to be written.
	encodedKey, encodedVal := encodeKeyValue(key, val)
	// The metadata object is a multimap, and previous values may exist, but for opentracing headers, we do not append
	// we just override.
	m[encodedKey] = []string{encodedVal}
}

// ForeachKey is a opentracing.TextMapReader interface that extracts values.
func (m metadataTextMap) ForeachKey(callback func(key, val string) error) error {
	for k, vv := range m {
		for _, v := range vv {
			if decodedKey, decodedVal, err := metadata.DecodeKeyValue(k, v); err == nil {
				if err = callback(decodedKey, decodedVal); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed decoding opentracing from gRPC metadata: %v", err)
			}
		}
	}
	return nil
}

// encodeKeyValue encodes key and value qualified for transmission via gRPC.
// note: copy pasted from private values of grpc.metadata
func encodeKeyValue(k, v string) (string, string) {
	k = strings.ToLower(k)
	if strings.HasSuffix(k, binHdrSuffix) {
		val := base64.StdEncoding.EncodeToString([]byte(v))
		v = string(val)
	}
	return k, v
}

func InitTrace(serviceName, reportAddr, sampleType string, rate float64) (err error) {

	transport, err := zipkin.NewHTTPTransport(
		reportAddr,
		zipkin.HTTPBatchSize(16),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init zipkin: %v\n", err)
		return
	}

	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  sampleType,
			Param: rate,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	r := jaeger.NewRemoteReporter(transport)
	tracer, closer, err := cfg.New(serviceName,
		config.Logger(jaeger.StdLogger),
		config.Reporter(r))
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init Jaeger: %v\n", err)
		return
	}

	_ = closer
	opentracing.SetGlobalTracer(tracer)
	return
}
