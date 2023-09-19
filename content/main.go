package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"content/handler"
	"content/subscriber"

	content "content/proto/content"
)

func main() {
	// New Service
	service := micro.NewService()

	// Initialise service
	service.Init()

	// Register Handler
	content.RegisterContentHandler(service.Server(), new(handler.Content))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("zusux.book.service.content", service.Server(), new(subscriber.Content))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
