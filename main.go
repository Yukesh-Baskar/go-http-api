package main

import (
	configs "github.com/http-crud/api/configs"
)

func main() {
	// `configs.LoadEnvVarsAndStartApp()` is a function call that loads environment variables and starts
	// the application. It is likely defined in the `configs` package and contains code to read
	// environment variables from a configuration file or the system environment, and then initializes and
	// starts the application with those values.
	configs.LoadEnvVarsAndStartApp()
}
