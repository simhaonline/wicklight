// +build windows
package main

import (
	"os"
	"os/signal"
	"syscall"
	"wicklight/quota"
)

func init() {
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT)

	go func() {
		<-sign
		quota.StoreQuota()
		os.Exit(0)
	}()
}
