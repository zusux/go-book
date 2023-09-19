package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"booklist/handler"
	"booklist/subscriber"

	booklist "booklist/proto/booklist"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	booklist.RegisterBooklistHandler(service.Server(), new(handler.Booklist))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("zusux.book.service.booklist", service.Server(), new(subscriber.Booklist))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
