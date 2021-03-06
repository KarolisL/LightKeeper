package main

import (
	"flag"
	"log"
	"os"

	"github.com/KarolisL/lightkeeper/pkg/daemon"
	"github.com/KarolisL/lightkeeper/pkg/daemon/config"
	"github.com/KarolisL/lightkeeper/pkg/plugins/input"
	_ "github.com/KarolisL/lightkeeper/pkg/plugins/input/file"
	"github.com/KarolisL/lightkeeper/pkg/plugins/output"
	_ "github.com/KarolisL/lightkeeper/pkg/plugins/output/stdout"
	_ "github.com/KarolisL/lightkeeper/pkg/plugins/output/telegram"
)

var configLocation string

func init() {
	flag.StringVar(&configLocation, "config", "", "path to config")
}

func main() {
	flag.Parse()
	if configLocation == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg, _ := config.NewConfigFromFile(configLocation)
	d, err := daemon.NewDaemon(cfg, input.Registry, output.Registry)
	if err != nil {
		log.Fatalf("Unable to start: %v", err)
	}

	d.Start()
	blockForever()
}

func blockForever() {
	<-make(chan struct{})
}
