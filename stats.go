package logging

import (
	"bytes"
	"expvar"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	StatsPrefix = "rps"
)

// By using StatsLog, you can print stats on stdout every second, which is sometimes handy to check the state
// of the server. The stats themselves are declared using the "expvar" package
// to use this function, just before starting your listeners, create a goroutine like this
// go logging.StatsLog()

func StatsLogInterval(seconds int, compact bool) {

	// If we are running in debug mode, do not clog the screen
	if IsDebug() {
		log.Println("disabling logger in debug mode")
		return
	}

	log.Println("starting logger")
	info := log.New(os.Stdout, "s:", log.Ldate|log.Ltime)

	sleepTime := time.Duration(seconds) * time.Second

	for _ = range time.Tick(sleepTime) {
		var buffer bytes.Buffer
		expvar.Do(func(k expvar.KeyValue) {
			// reset stats every nseconds
			shouldLog := !compact
			prev := k.Value.String()
			if v, ok := (k.Value).(*expvar.Int); ok {
				if prev != "0" {
					v.Set(0)
					shouldLog = true
				}
			}
			if strings.HasPrefix(k.Key, StatsPrefix) && shouldLog {
				buffer.WriteString(fmt.Sprintf("[%s %s] ", strings.TrimLeft(k.Key, StatsPrefix), prev))
			}
		})
		info.Println(buffer.String())
	}
}

func StatsLog() {
	StatsLogInterval(1, false)
}
