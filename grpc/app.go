package server

import (
	"context"

	"gitlab.com/promptech1/infuser-author/handler"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/model"
)

type appManagerServer struct {
	appHandler *handler.AppHandler
}

func (a appManagerServer) Create(ctx context.Context, req *grpc_author.AppReq) (*grpc_author.AppRes, error) {
	app := model.NewAppByGrpc(req)

	if err := a.appHandler.Create(app); err != nil {
		return &grpc_author.AppRes{Status: grpc_author.AppRes_ERROR}, nil
	}

	return &grpc_author.AppRes{Status: grpc_author.AppRes_OK}, nil
}

func (a appManagerServer) Update(ctx context.Context, req *grpc_author.AppReq) (*grpc_author.AppRes, error) {
	app := model.NewAppByGrpc(req)

	if err := a.appHandler.Update(app); err != nil {
		return &grpc_author.AppRes{Status: grpc_author.AppRes_ERROR}, nil
	}

	return &grpc_author.AppRes{Status: grpc_author.AppRes_OK}, nil
}

func (a appManagerServer) Destroy(ctx context.Context, req *grpc_author.AppReq) (*grpc_author.AppRes, error) {
	panic("implement me")
}

func newAppManagerServer(appHandler *handler.AppHandler) grpc_author.AppManagerServer {
	return &appManagerServer{appHandler: appHandler}
}
