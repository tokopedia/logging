package tracer

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/opentracing/opentracing-go"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
	"sourcegraph.com/sourcegraph/appdash/traceapp"

	"gopkg.in/tokopedia/logging.v1"
)

func InitAppdash(cfg *Config) {
	logging.Debug.Println("starting tracer on ", cfg.Port)
	go setupTracer(cfg.Port, cfg.TTL, cfg.ServerName)
}

// setupTracer must be called in a separate goroutine, as it's blocking
func setupTracer(appdashPort int, ttl int, server string) {

	// Tracer setup
	memStore := appdash.NewMemoryStore()

	// keep last hour of traces. In production, this will need to be modified.
	store := &appdash.RecentStore{
		MinEvictAge: time.Duration(ttl) * time.Second,
		DeleteStore: memStore,
	}
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		log.Fatalln("appdash", err)
	}

	collectorPort := l.Addr().String()
	logging.Debug.Println("collector listening on", collectorPort)

	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go cs.Start()

	if server == "" {
		server = fmt.Sprintf("http://localhost:%d", appdashPort)
	}

	appdashURL, err := url.Parse(server)
	tapp, err := traceapp.New(nil, appdashURL)
	if err != nil {
		log.Fatal(err)
	}
	tapp.Store = store
	tapp.Queryer = memStore

	tracer := appdashot.NewTracer(appdash.NewRemoteCollector(collectorPort))
	opentracing.InitGlobalTracer(tracer)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", appdashPort), tapp); err != nil {
		time.Sleep(15 * time.Second) // sleep 15 seconds so we don't run into port conflicts, but if we do again, exit
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appdashPort), tapp))
	}
}
