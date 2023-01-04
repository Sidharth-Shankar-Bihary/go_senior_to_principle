package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ramseyjiang/go_senior_to_principle/cmd/apiserver"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover from main(): [%+v]\n", r)
		}
	}()

	go apiserver.Main()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	fmt.Printf("main: shutting down server...")
}