package grpc

import (
	"context"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/service"
)

type userServer struct {
	userService service.UserService
}

func (u userServer) Login(ctx context.Context, req *grpc_author.UserReq) (*grpc_author.UserRes, error) {
	return u.userService.Login(req)
}

func (u userServer) Create(ctx context.Context, req *grpc_author.UserReq) (*grpc_author.UserRes, error) {
	return u.userService.Create(req)
}

func newUserServer(userService service.UserService) grpc_author.UserManagerServer {
	return &userServer{
		userService: userService,
	}
}
