package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kennydo/automatic-light-controller/cmd/daemon/app"
)

func main() {
	var config app.Config
	if _, err := toml.DecodeFile("configs/example.toml", &config); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Config: %+v", config)

}
