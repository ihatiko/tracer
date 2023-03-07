package tracer

import (
	"io"

	opentracingClient "github.com/opentracing/opentracing-go"
	jaegerClient "github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

type Tracer struct {
	TracerClient opentracingClient.Tracer
	Closer       io.Closer
}

func (cfg *Config) NewTracer(serviceName string) (*Tracer, error) {
	jaegerCfgInstance := jaegerCfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaegerClient.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:           cfg.LogSpans,
			LocalAgentHostPort: cfg.Host,
		},
	}
	tracer, ioCloser, err := jaegerCfgInstance.NewTracer(
		jaegerCfg.Logger(jaegerLog.StdLogger),
		jaegerCfg.Metrics(metrics.NullFactory),
	)
	opentracingClient.SetGlobalTracer(tracer)
	return &Tracer{
		TracerClient: tracer,
		Closer:       ioCloser,
	}, err
}
