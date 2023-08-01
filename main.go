package main

import (
	"flag"
	"github.com/gogoclouds/go-web/intermal/boot"
)

var config = flag.String("config", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	boot.Run(*config)
}