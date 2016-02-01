package logging

import (
	"bytes"
	"expvar"
	"fmt"
	"log"
	"strings"
	"time"
  "os"
)

// By using StatsLog, you can print stats on stdout every second, which is sometimes handy to check the state
// of the server. The stats themselves are declared using the "expvar" package
// to use this function, just before starting your listeners, create a goroutine like this
// go logging.StatsLog()
func StatsLog() {

  // If we are running in debug mode, do not clog the screen
	if IsDebug() {
    log.Println("disabling logger in debug mode")
		return
	}

  log.Println("starting logger")
  info := log.New(os.Stdout, "s:", log.Ldate|log.Ltime)


	for _ = range time.Tick(time.Second) {
		var buffer bytes.Buffer
		expvar.Do(func(k expvar.KeyValue) {
			if strings.HasPrefix(k.Key, "rps") {
				buffer.WriteString(fmt.Sprintf("[%s %s] ", k.Key, k.Value))
				// reset stats every second
				if v, ok := (k.Value).(*expvar.Int); ok {
					v.Set(0)
				}
			}
		})
		info.Println(buffer.String())
	}

}
