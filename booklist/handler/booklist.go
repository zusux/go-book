package handler

import (
	"context"
	"math"

	log "github.com/micro/go-micro/v2/logger"
	"booklist/config"

	booklist "booklist/proto/booklist"
)

type Booklist struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Booklist) Call(ctx context.Context, req *booklist.BookRequest, rsp *booklist.Books) error {
	log.Info("Received Booklist.Call request")
	log.Info(req)
	db := config.GetDb().Debug()
	tdb := db.Table("b_book b").Select("b.*,c.cat_name").Where("b.status = ?",1).
		Joins("left join b_cat c on b.cat_id = c.id")
	if req.Id > 0{
		tdb = tdb.Where("b.id = ?", req.Id)
	}
	if req.CatId >0 {
		tdb = tdb.Where("cat_id = ?", req.CatId)
	}
	if req.Name != ""{
		tdb = tdb.Where("name = ?", req.Name)
	}
	if req.Author != ""{
		tdb = tdb.Where("author = ?", req.Name)
	}
	if req.IsHot >0{
		tdb = tdb.Where("is_hot = ?", req.IsHot)
	}
	if req.IsNew >0{
		tdb = tdb.Where("is_new = ?", req.IsNew)
	}
	if req.IsOver >0{
		tdb = tdb.Where("is_over = ?", req.IsOver)
	}
	if req.Limit > 0{
		tdb = tdb.Limit(req.Limit)
	}
	if req.Page >= 0{
		tdb = tdb.Offset(req.Page * int32(math.Abs(float64(req.Limit))))
	}

	if req.Order != ""{
		tdb = tdb.Order(req.Order)
	}

	tdb.Find(&rsp.BookResponses)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Booklist) Stream(ctx context.Context, req *booklist.StreamingRequest, stream booklist.Booklist_StreamStream) error {
	log.Infof("Received Booklist.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&booklist.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Booklist) PingPong(ctx context.Context, stream booklist.Booklist_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&booklist.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
