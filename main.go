package main

import (
	"os"
	"os/signal"
	"syscall"
	"wicklight/config"
	"wicklight/quota"
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
	quota.InitQuota()
	defer quota.StoreQuota()
	server.Run()
}

func init() {
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT)

	go func() {
		<-sign
		quota.StoreQuota()
		os.Exit(0)
	}()
}
