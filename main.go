package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
)

func main() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		systray.Quit()
	}()

	systray.Run(onReady, onExit)
}
