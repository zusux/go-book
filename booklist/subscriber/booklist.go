package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	booklist "booklist/proto/booklist"
)

type Booklist struct{}

func (e *Booklist) Handle(ctx context.Context, msg *booklist.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *booklist.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
