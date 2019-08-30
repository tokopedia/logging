// Package logging provides common functionality as log rotation, conditional debug logging etc.
// To initialize this package, just import it as
// import "gopkg.in/tokopedia/logging.v1
package logging

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var stdoutLog string
var stderrLog string

// separating debug log into a specific file to ease monitoring process
var debugLog string

var debugFlag bool
var versionFlag bool

// Debug global logger for debug messages
// logging.Debug.Println("debug message")
// debug messages are printed only when the program is started with -debug flag
var Debug *log.Logger

// Init installs the command line options for setting output and error log paths, and exposes
// logging.Debug, which can be used to add code for debug
func init() {
	flag.StringVar(&stdoutLog, "l", "", "log file for stdout")
	flag.StringVar(&stderrLog, "e", "", "log file for stderr")
	flag.StringVar(&debugLog, "d", "", "log file for debug") // only if to use debugLog
	flag.BoolVar(&versionFlag, "version", false, "binary version")
	flag.BoolVar(&debugFlag, "debug", false, "enable debug logging")

	Debug = log.New(ioutil.Discard, "", 0)

	// if running with socketmaster, reload is really not needed
	if fd := os.Getenv("EINHORN_FDS"); fd == "" {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP) // listen for sighup
		go sigHandler(c)
	}
}

func sigHandler(c chan os.Signal) {
	// Block until a signal is received.
	for s := range c {
		log.Println("Reloading on :", s)
		LogInit()
	}
}

// LogInit App must call LogInit once to setup log redirection
func LogInit() {

	if versionFlag {
		fmt.Println(appVersion())
		os.Exit(0)
	}

	if stdoutLog != stderrLog && stdoutLog != "" {
		log.Println("Log Init: using ", stdoutLog, stderrLog)
	}

	_, _ = reopen(1, stdoutLog)
	_, _ = reopen(2, stderrLog)

	SetDebug(debugFlag)
}

// SetDebug set debug
func SetDebug(enabled bool) {
	if enabled {
		debugFlag = true
		f, err := reopen(0, debugLog)
		if f == nil || err != nil {
			f = os.Stdout
		}

		Debug = log.New(f, "debug:", log.Ldate|log.Ltime|log.Lshortfile)
		Debug.Println("---- debug mode ----")
	}
}

// IsDebug Determine if we are running in debug mode or not
func IsDebug() bool {
	return debugFlag
}

func reopen(fd int, filename string) (*os.File, error) {
	if filename == "" {
		return nil, fmt.Errorf("Empty log file for fd: %d", fd)
	}

	logFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if fd != 0 && err != nil {
		// do not terminate in case of debug file open error
		log.Println("Error in opening ", filename, err)
		os.Exit(2)
	}

	if fd != 0 {
		if err = syscall.Dup2(int(logFile.Fd()), fd); err != nil {
			log.Println("Failed to dup", filename)
		}
	}

	return logFile, err
}
