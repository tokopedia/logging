package tracer

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func InitJaeger(tcfg *Config) {
	// Initialize tracer with a logger and a metrics factory, returning closer
	appname := path.Base(os.Args[0])
	env := os.Getenv("TKPENV")
	if env == "" {
		env = "development"
	}

	// Default sampler is const, we sample everything, do not use in production
	sampler := jaegercfg.SamplerConfig{
		Type:  jaeger.SamplerTypeConst,
		Param: 1,
	}

	if additionalCfg := tcfg.Extras; additionalCfg != "" {
		fields := strings.Split(additionalCfg, ":")
		switch fields[0] {
		case jaeger.SamplerTypeConst:
			sampler.Type = jaeger.SamplerTypeConst
		case jaeger.SamplerTypeRateLimiting:
			sampler.Type = jaeger.SamplerTypeRateLimiting
		case jaeger.SamplerTypeProbabilistic:
			sampler.Type = jaeger.SamplerTypeProbabilistic
		case jaeger.SamplerTypeLowerBound:
			sampler.Type = jaeger.SamplerTypeLowerBound
		default:
			sampler.Type = jaeger.SamplerTypeConst
		}

		if param, err := strconv.ParseFloat(fields[1], 64); err == nil {
			sampler.Param = param
		}
	}

	cfg := jaegercfg.Configuration{
		Sampler: &sampler,
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jMetricsFactory := metrics.NullFactory
	if tcfg.Verbose {
		log.Println("jaeger:", appname, env, sampler.Type, sampler.Param)
		jaegercfg.Logger(jaegerlog.StdLogger)
	} else {
		jaegercfg.Logger(jaegerlog.NullLogger)
	}

	_, err := cfg.InitGlobalTracer(
		appname+"-"+env,
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	//defer closer.Close()
}
