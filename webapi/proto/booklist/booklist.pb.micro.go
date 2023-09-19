// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: booklist.proto

package zusux_book_service_booklist

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

// Api Endpoints for Booklist service

func NewBooklistEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Booklist service

type BooklistService interface {
	Call(ctx context.Context, in *BookRequest, opts ...client.CallOption) (*Books, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Booklist_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Booklist_PingPongService, error)
}

type booklistService struct {
	c    client.Client
	name string
}

func NewBooklistService(name string, c client.Client) BooklistService {
	return &booklistService{
		c:    c,
		name: name,
	}
}

func (c *booklistService) Call(ctx context.Context, in *BookRequest, opts ...client.CallOption) (*Books, error) {
	req := c.c.NewRequest(c.name, "Booklist.Call", in)
	out := new(Books)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *booklistService) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Booklist_StreamService, error) {
	req := c.c.NewRequest(c.name, "Booklist.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &booklistServiceStream{stream}, nil
}

type Booklist_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type booklistServiceStream struct {
	stream client.Stream
}

func (x *booklistServiceStream) Close() error {
	return x.stream.Close()
}

func (x *booklistServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *booklistServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *booklistServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *booklistServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *booklistService) PingPong(ctx context.Context, opts ...client.CallOption) (Booklist_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Booklist.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &booklistServicePingPong{stream}, nil
}

type Booklist_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type booklistServicePingPong struct {
	stream client.Stream
}

func (x *booklistServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *booklistServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *booklistServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *booklistServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *booklistServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *booklistServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Booklist service

type BooklistHandler interface {
	Call(context.Context, *BookRequest, *Books) error
	Stream(context.Context, *StreamingRequest, Booklist_StreamStream) error
	PingPong(context.Context, Booklist_PingPongStream) error
}

func RegisterBooklistHandler(s server.Server, hdlr BooklistHandler, opts ...server.HandlerOption) error {
	type booklist interface {
		Call(ctx context.Context, in *BookRequest, out *Books) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Booklist struct {
		booklist
	}
	h := &booklistHandler{hdlr}
	return s.Handle(s.NewHandler(&Booklist{h}, opts...))
}

type booklistHandler struct {
	BooklistHandler
}

func (h *booklistHandler) Call(ctx context.Context, in *BookRequest, out *Books) error {
	return h.BooklistHandler.Call(ctx, in, out)
}

func (h *booklistHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.BooklistHandler.Stream(ctx, m, &booklistStreamStream{stream})
}

type Booklist_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type booklistStreamStream struct {
	stream server.Stream
}

func (x *booklistStreamStream) Close() error {
	return x.stream.Close()
}

func (x *booklistStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *booklistStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *booklistStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *booklistStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *booklistHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.BooklistHandler.PingPong(ctx, &booklistPingPongStream{stream})
}

type Booklist_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type booklistPingPongStream struct {
	stream server.Stream
}

func (x *booklistPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *booklistPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *booklistPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *booklistPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *booklistPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *booklistPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}