package server

import (
	"context"
	"gitlab.com/promptech1/infuser-author/handler"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/model"
)

type apiAuthServer struct {
	handler *handler.AppTokenHandler
}

func (a *apiAuthServer) Auth(ctx context.Context, req *grpc_author.ApiAuthReq) (*grpc_author.ApiAuthRes, error) {
	token := model.Token{Token: req.Token, IsDel: false}
	operation := model.Operation{
		EndPoint: req.OperationUrl, IsDel: false,
		App: model.App{NameSpace: req.NameSpace, IsDel: false},
	}

	authCode := a.handler.CheckAppToken(&token, &operation)

	res := &grpc_author.ApiAuthRes{
		Code: authCode,
	}

	return res, nil
}

func newApiAuthServer(handler *handler.AppTokenHandler) grpc_author.ApiAuthServiceServer {
	return &apiAuthServer{
		handler: handler,
	}
}
