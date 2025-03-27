package main

import (
	"duo_elements/internal/app"
	"flag"
)

func main() {
	configFile := flag.String("config", ".env", "Path to configuration file")
	flag.Parse()

	app.Run(*configFile)
}
