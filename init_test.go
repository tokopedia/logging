package logging

func ExampleLogInit() {
  logging.LogInit()
  // Setup log redirection, debug logger can be called now
  logging.Debug.Println("this will only show if app was started with -debug")
}
