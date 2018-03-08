package logging

import (
	"log"
	"os"

	gcfg "gopkg.in/gcfg.v1"
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

func MustReadModuleConfig(cfg interface{}, paths []string, module string) {
	res := false
	for _, path := range paths {
		res = ReadModuleConfig(cfg, path, module)
		if res == true {
			break
		}
	}

	if res == false {
		log.Fatalln("couldn't read config for ", os.Getenv("TKPENV"))
	}
}
