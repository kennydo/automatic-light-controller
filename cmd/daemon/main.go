package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kennydo/automatic-light-controller/cmd/daemon/app"
	"github.com/kennydo/automatic-light-controller/lib/huebridge"
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

	fmt.Printf("Config: %+v", config)

	hueBridge, err := huebridge.New(config.HueBridge.IPAddress.String(), config.HueBridge.Username)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bridge: %v,", hueBridge)
}
