package tracer

type Config struct {
	Enabled      bool   // set to disabled to not use any tracer
	Name         string // which tracer to use, appdash or jaeger
	Port         int    // tracer listen port
	TTL          int    // seconds to live for the trace
	ExcludeRegex string // url's matching this will be skipped.
	ServerName   string // defaults to the Name specified in  ServerConfig
	Verbose      bool   // enable detailed logging on tracer when set
	Extras       string // extra config specific to tracer (e.g. appdash or jaeger)
}

func Init(tracerCfg *Config) {
	if tracerCfg.Enabled {
		// default is appdash
		if tracerCfg.Name == "" || tracerCfg.Name == "appdash" {
			InitAppdash(tracerCfg)
		} else {
			InitJaeger(tracerCfg)
		}
	}
}
