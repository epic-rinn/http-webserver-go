package main

import (
	"log"
	"net"
	"os"
)

var _ = net.Listen
var _ = os.Exit

type application struct {
	config *Config
}

func main() {
	var cfg Config

	LoadEnvs(&cfg)

	app := &application{
		config: &cfg,
	}

	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
