package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.com/promptech1/infuser-author/handler"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/model"
)

type appManagerServer struct {
	appHandler *handler.AppHandler
}

func newAppManagerServer(appHandler *handler.AppHandler) grpc_author.AppManagerServer {
	return &appManagerServer{appHandler: appHandler}
}

func (a appManagerServer) Create(ctx context.Context, req *grpc_author.AppReq) (*grpc_author.AppRes, error) {
	app := model.NewAppByGrpc(req)

	if err := a.appHandler.Create(app); err != nil {
		a.appHandler.Ctx.Logger.WithFields(logrus.Fields{
			"module":   "appManagerServer",
			"function": "Create",
		}).Info(err)
		return &grpc_author.AppRes{Status: grpc_author.AppRes_ERROR}, nil
	}

	return &grpc_author.AppRes{Status: grpc_author.AppRes_OK}, nil
}

func (a appManagerServer) Update(ctx context.Context, req *grpc_author.AppReq) (*grpc_author.AppRes, error) {
	app := model.NewAppByGrpc(req)

	if err := a.appHandler.Update(app); err != nil {
		a.appHandler.Ctx.Logger.WithFields(logrus.Fields{
			"module":   "appManagerServer",
			"function": "Update",
		}).Info(err)
		return &grpc_author.AppRes{Status: grpc_author.AppRes_ERROR}, nil
	}

	return &grpc_author.AppRes{Status: grpc_author.AppRes_OK}, nil
}

func (a appManagerServer) Destroy(ctx context.Context, req *grpc_author.AppReq) (*grpc_author.AppRes, error) {
	if err := a.appHandler.Destroy(uint(req.AppId)); err != nil {
		a.appHandler.Ctx.Logger.WithFields(logrus.Fields{
			"module":   "appManagerServer",
			"function": "Destroy",
		}).Info(err)
		return &grpc_author.AppRes{Status: grpc_author.AppRes_ERROR}, nil
	}

	return &grpc_author.AppRes{Status: grpc_author.AppRes_OK}, nil
}
