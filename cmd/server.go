package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"webinar/graphql/server/internal"
)

func main() {
	service := &internal.Service{}
	if err := service.Setup(); err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		if err := service.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	if err := service.Shutdown(); err != nil {
		log.Fatal(err.Error())
	}
}
