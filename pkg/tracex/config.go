package tracex

import (
	"fmt"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"log"
)

type Config struct {
	Addr        string `yaml:"addr"`
	ServiceName string `yaml:"service_name"`
	Rate        int    `yaml:"rate"` //0-100
}

func WithConfig(conf Config) Option {
	return func(t *Tracer) {
		t.conf = conf
	}
}

type Option func(t *Tracer)

func Must(opts ...Option) *Tracer {
	tracer := newTracer()
	for _, opt := range opts {
		opt(tracer)
	}
	err := tracer.configuration()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	return tracer
}

func (t *Tracer) configuration() error {
	tracer, closer, err := config.Configuration{
		ServiceName: t.conf.ServiceName,
		Sampler: &config.SamplerConfig{
			Type: "probabilistic", Param: float64(t.conf.Rate) / 100.0,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: t.conf.Addr,
			LogSpans:           false,
		},
	}.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return err
	}
	t.Closer = closer
	t.Tracer = tracer
	return nil
}
