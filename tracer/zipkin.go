package tracer

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/go-kit/kit/log"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

// References
//https://golang.hotexamples.com/site/file?hash=0xa4dcb4a0a61c426467f5b4b50143174aa393f0f74977c9fd440e237da6d28e6f&fullName=examples/main.go&project=zyanho/kit
// https://github.com/openzipkin/zipkin-go/blob/master/example_httpserver_test.go
// https://github.com/openzipkin-contrib/zipkin-go-opentracing/blob/master/examples/cli_with_2_services/svc2/implementation.go
//https://medium.com/@Ankitthakur/apache-kafka-installation-on-mac-using-homebrew-a367cdefd273 Kafka setup

// setup
// brew cask install homebrew/cask-versions/adoptopenjdk8
// brew install kafka
// docker run -d -p 9411:9411 openzipkin/zipkinm and go to http://localhost:9411/

const debug = false

const zipkinType = "api"

// zipkinHTTPEndpoint is default zipkin endpoint
const zipkinHTTPEndpoint = "http://localhost:9411/api/v1/spans"

// const zipkinAddr = "127.0.0.1:9092"

func InitZipkin(tcfg *Config) {
	// Initialize tracer with a logger and a metrics factory, returning closer

	logText := "Initialize Zipkin"

	var httpCollector, envName string
	if additionalCfg := tcfg.Extras; additionalCfg != "" {
		//Extras = "api:ENV"
		extras := strings.Split(additionalCfg, ":")
		httpCollector = extras[0]

		if len(extras) > 0 {
			envName = extras[1]
		}
	}

	appname := path.Base(os.Args[0])
	env := os.Getenv(envName)
	if env == "" {
		env = "development"
	}

	var collector zipkin.Collector
	var err error

	if httpCollector == zipkinType {
		collector, err = zipkin.NewHTTPCollector(zipkinHTTPEndpoint)
		logText = logText + " API"
	} else {
		collector, err = zipkin.NewKafkaCollector(
			strings.Split("127.0.0.1:9092", ","), // default kafka port
			zipkin.KafkaLogger(log.NewNopLogger()),
		)
		logText = logText + " Kafka"
	}

	if err != nil {
		fmt.Println("err", "unable to create collector", "fatal", err)
		return
	}

	fmt.Println(logText)
	var hostPort string
	if tcfg.ServerName == "" {
		hostPort = fmt.Sprintf("127.0.0.1:%d", tcfg.Port)
	}

	serviceName := env + "-" + appname

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)

	// Create our tracer.
	zipkinTracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(true), // always  using same span same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)
		zipkin.TraceID128Bit(true),
	)
	if err != nil {
		fmt.Printf("Could not initialise Zipkin tracer: %+v\n", err)
		return
	}

	// Explicitly set our tracer to be the default tracer.
	opentracing.InitGlobalTracer(zipkinTracer)
	//collector.Close()
}
