package main

import (
	"flag"
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
	_, err := toml.DecodeFile(configPath, config)//

	if err != nil {
		log.Fatal(err)
	}

	if err := apisrever.Start(config); err != nil {
		log.Fatal(err)
	}
}
