package handler

import (
	"content/config"
	"context"
	log "github.com/micro/go-micro/v2/logger"

	content "content/proto/content"
)

type Content struct{
	Id int32 `gorm:"column:id"`
	ChapterId int32 `gorm:"column:chapter_id"`
	BookId int32 `gorm:"column:book_id"`
	Cont string `gorm:"column:content"`
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Content) Call(ctx context.Context, req *content.Request, rsp *content.Response) error {
	//log.Info("Received Content.Call request")
	//rsp.Msg = "Hello " + req.Name
	log.Info(req)
	db := config.GetDb().Table("b_content")
	if req.BookId >0 {
		db = db.Where("book_id = ?", req.BookId)
	}
	if req.ChapterId > 0{
		db = db.Where("chapter_id = ?", req.ChapterId)
	}
	db.First(&rsp)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Content) Stream(ctx context.Context, req *content.StreamingRequest, stream content.Content_StreamStream) error {
	log.Infof("Received Content.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&content.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Content) PingPong(ctx context.Context, stream content.Content_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&content.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
