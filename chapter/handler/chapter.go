package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"
	"chapter/config"

	chapter "chapter/proto/chapter"
)

type Chapter struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Chapter) Call(ctx context.Context, req *chapter.ChapterRequest, rsp *chapter.Chapters) error {
	log.Info("Received Chapter.Call request")
	log.Info("req",req.Id,req.BookId)
	db := config.GetDb().Table("b_chapter").Debug()
	if req.BookId >0 {
		db = db.Where("book_id = ?", req.BookId)
	}
	if req.Id > 0{
		db = db.Where("id = ?", req.Id)
	}
	db.Order("sort").Find(&rsp.ChapterResponses)

	log.Info("rsp",rsp)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Chapter) Stream(ctx context.Context, req *chapter.StreamingRequest, stream chapter.Chapter_StreamStream) error {
	log.Infof("Received Chapter.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&chapter.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Chapter) PingPong(ctx context.Context, stream chapter.Chapter_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&chapter.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
