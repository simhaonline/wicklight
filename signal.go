// +build linux darwin freebsd netbsd openbsd
package main

import (
	"os"
	"os/signal"
	"syscall"
	"wicklight/quota"
)

func init() {
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sign
		quota.StoreQuota()
		os.Exit(0)
	}()
}
