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
	app := new(model.App)
	app.Id = uint(req.AppId)
	app.NameSpace = req.NameSpace

	if len(req.Operations) > 0 {
		for _, operation := range req.Operations {
			app.Operations = append(app.Operations, model.Operation{
				Id:       uint(operation.OperationId),
				AppId:    app.Id,
				EndPoint: operation.EndPoint,
			})
		}
	}

	if len(req.Traffics) > 0 {
		for _, traffic := range req.Traffics {
			app.Traffics = append(app.Traffics, model.Traffic{
				AppId: app.Id,
				Unit:  traffic.Unit,
				Val:   uint(traffic.Value),
				Seq:   uint(traffic.Seq),
			})
		}
	}

	if err := a.appHandler.Create(app); err != nil {
		return &grpc_author.AppRes{Status: grpc_author.AppRes_ERROR}, nil
	}

	return &grpc_author.AppRes{Status: grpc_author.AppRes_OK}, nil
}

func (a appManagerServer) Update(ctx context.Context, req *grpc_author.AppReq) (*grpc_author.AppRes, error) {
	panic("implement me")
}

func (a appManagerServer) Destroy(ctx context.Context, req *grpc_author.AppReq) (*grpc_author.AppRes, error) {
	panic("implement me")
}

func newAppManagerServer(appHandler *handler.AppHandler) grpc_author.AppManagerServer {
	return &appManagerServer{appHandler: appHandler}
}
