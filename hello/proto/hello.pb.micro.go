// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: hello.proto

package hello

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/asim/go-micro/v3/api"
	client "github.com/asim/go-micro/v3/client"
	server "github.com/asim/go-micro/v3/server"
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

// Api Endpoints for Meet service

func NewMeetEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Meet service

type MeetService interface {
	Hello(ctx context.Context, in *User, opts ...client.CallOption) (*User, error)
}

type meetService struct {
	c    client.Client
	name string
}

func NewMeetService(name string, c client.Client) MeetService {
	return &meetService{
		c:    c,
		name: name,
	}
}

func (c *meetService) Hello(ctx context.Context, in *User, opts ...client.CallOption) (*User, error) {
	req := c.c.NewRequest(c.name, "Meet.Hello", in)
	out := new(User)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Meet service

type MeetHandler interface {
	Hello(context.Context, *User, *User) error
}

func RegisterMeetHandler(s server.Server, hdlr MeetHandler, opts ...server.HandlerOption) error {
	type meet interface {
		Hello(ctx context.Context, in *User, out *User) error
	}
	type Meet struct {
		meet
	}
	h := &meetHandler{hdlr}
	return s.Handle(s.NewHandler(&Meet{h}, opts...))
}

type meetHandler struct {
	MeetHandler
}

func (h *meetHandler) Hello(ctx context.Context, in *User, out *User) error {
	return h.MeetHandler.Hello(ctx, in, out)
}
