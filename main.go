package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chetan177/pismo/api"
)

func main() {
	log.Println("server is running")
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	server := api.NewApiServer()
	server.Start()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	<-done

	server.Stop()
}
