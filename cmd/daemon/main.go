package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/kennydo/automatic-light-controller/cmd/daemon/app"
	"log"
)

func main() {
	var config app.Config
	if _, err := toml.DecodeFile("configs/example.toml", &config); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Config: %+v", config)

}
