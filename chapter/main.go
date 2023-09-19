package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"chapter/handler"
	"chapter/subscriber"

	chapter "chapter/proto/chapter"
)

func main() {
	// New Service
	service := micro.NewService()

	// Initialise service
	service.Init()

	// Register Handler
	chapter.RegisterChapterHandler(service.Server(), new(handler.Chapter))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("zusux.book.service.chapter", service.Server(), new(subscriber.Chapter))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
