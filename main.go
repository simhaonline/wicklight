package main

import (
	"os"
	"wicklight/config"
	"wicklight/server"
	"wicklight/version"
)

func main() {
	if len(os.Args) != 2 {
		println(version.Print())
		println(version.Help())
		os.Exit(0)
	}
	config.ReadConfig(os.Args[1])
	server.Run()
}
