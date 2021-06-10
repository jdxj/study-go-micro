package handler

import (
	"context"
	"fmt"

	proto "study_go_micro/hello/proto"
)

type Helloer struct{}

func (h *Helloer) Hello(ctx context.Context, req *proto.User, reply *proto.User) error {
	reply.Name = fmt.Sprintf("hello %s", req.Name)
	return nil
}
