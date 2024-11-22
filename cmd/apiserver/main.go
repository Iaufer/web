package main

import (
	"flag"
	"fmt"
	"log"
	"web/internal/app/apisrever"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apisrever.NewConfig()
	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Start...")

	if err := apisrever.Start(config); err != nil {
		log.Fatal(err)
	}
}
