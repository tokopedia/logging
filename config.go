package logging

import (
	gcfg "gopkg.in/gcfg.v1"
	"os"
)

func ReadModuleConfig(cfg interface{}, path string, module string) bool {
	environ := os.Getenv("TKPENV")
	if environ == "" {
		environ = "development"
	}

	debug := Debug.Println

	fname := path + "/" + module + "." + environ + ".ini"
	err := gcfg.ReadFileInto(cfg, fname)
	if err == nil {
		debug("read config from ", fname)
		return true
	}
	debug(err)
	return false
}
