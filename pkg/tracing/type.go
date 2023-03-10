package tracing

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func New() (io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: viper.GetString("service.name"),
		RPCMetrics:  true,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", viper.GetString("tracing.host"), viper.GetUint32("tracing.port")),
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new tracer")
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
