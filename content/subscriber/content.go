package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	content "content/proto/content"
)

type Content struct{}

func (e *Content) Handle(ctx context.Context, msg *content.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *content.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
