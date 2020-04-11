package main

import (
	"flag"
	"fmt"
	"github.com/KarolisL/lightkeeper/pkg/daemon"
	"os"
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

	cfg, _ := daemon.NewConfigFromFile(configLocation)

	fmt.Printf("%+v", cfg)

}
