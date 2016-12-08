package logging

import (
	"testing"
)

func ExampleLogInit() {
	LogInit()
	// Setup log redirection, debug logger can be called now
	Debug.Println("this will only show if app was started with -debug")
}

func TestLogInit(t *testing.T) {
	Debug.Println("this will only show if app was started with -debug")
	LogInit()
	// Setup log redirection, debug logger can be called now
	Debug.Println("this will only show if app was started with -debug")
}
