package main

import (
	"os"

	"github.com/imeplusplus/dont-panic-api/app"
)

func main() {

	app := app.App{}
	app.Initialize()

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	app.Run(listenAddr)
}
