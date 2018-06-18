package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kennydo/automatic-light-controller/cmd/daemon/app"
)

var configPath = flag.String("config", "", "Path to a config file")

func main() {
	flag.Parse()

	if *configPath == "" {
		log.Fatal("Config file must be provided")
	}

	var config app.Config
	if _, err := toml.DecodeFile(*configPath, &config); err != nil {
		log.Fatal(err)
	}

	app, err := app.New(&config)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
