// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: chapter.proto

package zusux_book_service_chapter

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Chapter service

func NewChapterEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Chapter service

type ChapterService interface {
	Call(ctx context.Context, in *ChapterRequest, opts ...client.CallOption) (*Chapters, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Chapter_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Chapter_PingPongService, error)
}

type chapterService struct {
	c    client.Client
	name string
}

func NewChapterService(name string, c client.Client) ChapterService {
	return &chapterService{
		c:    c,
		name: name,
	}
}

func (c *chapterService) Call(ctx context.Context, in *ChapterRequest, opts ...client.CallOption) (*Chapters, error) {
	req := c.c.NewRequest(c.name, "Chapter.Call", in)
	out := new(Chapters)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chapterService) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Chapter_StreamService, error) {
	req := c.c.NewRequest(c.name, "Chapter.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &chapterServiceStream{stream}, nil
}

type Chapter_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type chapterServiceStream struct {
	stream client.Stream
}

func (x *chapterServiceStream) Close() error {
	return x.stream.Close()
}

func (x *chapterServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *chapterServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *chapterServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *chapterServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chapterService) PingPong(ctx context.Context, opts ...client.CallOption) (Chapter_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Chapter.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &chapterServicePingPong{stream}, nil
}

type Chapter_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type chapterServicePingPong struct {
	stream client.Stream
}

func (x *chapterServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *chapterServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *chapterServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *chapterServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *chapterServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *chapterServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Chapter service

type ChapterHandler interface {
	Call(context.Context, *ChapterRequest, *Chapters) error
	Stream(context.Context, *StreamingRequest, Chapter_StreamStream) error
	PingPong(context.Context, Chapter_PingPongStream) error
}

func RegisterChapterHandler(s server.Server, hdlr ChapterHandler, opts ...server.HandlerOption) error {
	type chapter interface {
		Call(ctx context.Context, in *ChapterRequest, out *Chapters) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Chapter struct {
		chapter
	}
	h := &chapterHandler{hdlr}
	return s.Handle(s.NewHandler(&Chapter{h}, opts...))
}

type chapterHandler struct {
	ChapterHandler
}

func (h *chapterHandler) Call(ctx context.Context, in *ChapterRequest, out *Chapters) error {
	return h.ChapterHandler.Call(ctx, in, out)
}

func (h *chapterHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.ChapterHandler.Stream(ctx, m, &chapterStreamStream{stream})
}

type Chapter_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type chapterStreamStream struct {
	stream server.Stream
}

func (x *chapterStreamStream) Close() error {
	return x.stream.Close()
}

func (x *chapterStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *chapterStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *chapterStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *chapterStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *chapterHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.ChapterHandler.PingPong(ctx, &chapterPingPongStream{stream})
}

type Chapter_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type chapterPingPongStream struct {
	stream server.Stream
}

func (x *chapterPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *chapterPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *chapterPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *chapterPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *chapterPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *chapterPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}