package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	chapter "chapter/proto/chapter"
)

type Chapter struct{}

func (e *Chapter) Handle(ctx context.Context, msg *chapter.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *chapter.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
