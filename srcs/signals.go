package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var shutdownSignal = make(chan struct{})

func signalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Printf("\nReceived shutdown signal, exiting...\n")
	close(shutdownSignal)
}
