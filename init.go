// This package handles log rotation
package logging

import (
  "os"
  "syscall"
  "flag"
  "log"
  "io/ioutil"
  "os/signal"
)

var stdoutLog string
var stderrLog string
var debugFlag bool
var Debug *log.Logger

func Init() {
  flag.StringVar(&stdoutLog,"l","","log file for stdout")
  flag.StringVar(&stderrLog,"e","","log file for stderr")
  flag.BoolVar(&debugFlag,"debug",false,"enable debug logging")

  c := make(chan os.Signal, 1)
  signal.Notify(c, syscall.SIGHUP) // listen for sighup
  go sigHandler(c)
}

func sigHandler(c chan os.Signal) {
  // Block until a signal is received.
  for s := range c {
    log.Println("Reloading on :", s)
    LogInit()
  }
}

func LogInit() {
  log.Println("Log Init: using ",stdoutLog,stderrLog)
  reopen(1,stdoutLog)
  reopen(2,stderrLog)

  if !debugFlag {
    Debug = log.New(ioutil.Discard,"",0)
  } else {
    Debug = log.New(os.Stdout,"debug:",log.Ldate|log.Ltime)
    Debug.Println("---- debug mode ----")
  }
}

func reopen(fd int,filename string) {
  if filename == "" {
    return
  }

  logFile,err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)

  if (err != nil) {
    log.Println("Error in opening ",filename,err)
    os.Exit(2)
  }

  syscall.Dup2(int(logFile.Fd()), fd)
}
