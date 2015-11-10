package logging

import (
  gcfg "gopkg.in/gcfg.v1"
  "os"
  "log"
)

func ReadModuleConfig(cfg interface{}, path string, module string) bool {
  environ := os.Getenv("TKPENV")
  if environ == "" {
    environ = "development"
  }

  fname := path + "/" + module + "." + environ + ".ini"
  err := gcfg.ReadFileInto(cfg,fname)
  if err == nil {
    log.Println("read config from ", fname)
    return true
  }
  log.Println(err)
  return false
}
