package main

import (
	"flag"
	"github.com/gogoclouds/go-web/intermal/app"
)

var config = flag.String("config", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	app.Run(*config)
}