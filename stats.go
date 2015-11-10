package logging

import (
	"bytes"
	"expvar"
	"fmt"
	"log"
	"strings"
	"time"
)

func StatsLog() {

	if IsDebug() {
		return
	}

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
		log.Println(buffer.String())
	}

}
