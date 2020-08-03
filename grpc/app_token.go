package grpc

import (
	"context"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/service"
)

type appTokenServer struct {
	appTokenService service.AppTokenService
}

func newAppTokenServer(appTokenService service.AppTokenService) grpc_author.AppTokenManagerServer {
	return &appTokenServer{appTokenService: appTokenService}
}

func (a appTokenServer) Regist(ctx context.Context, req *grpc_author.AppTokenReq) (*grpc_author.AppTokenRes, error) {
	a.appTokenService.Regist(req.Token, req.NameSpace)

	return &grpc_author.AppTokenRes{Code: 200}, nil
}

func (a appTokenServer) Destroy(ctx context.Context, req *grpc_author.AppTokenReq) (*grpc_author.AppTokenRes, error) {
	panic("implement me")
}
