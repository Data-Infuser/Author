package server

import (
	"context"
	"gitlab.com/promptech1/infuser-author/handler"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
)

type apiAuthServer struct {
	handler *handler.AppTokenHandler
}

func (a *apiAuthServer) Auth(ctx context.Context, req *grpc_author.ApiAuthReq) (*grpc_author.ApiAuthRes, error) {
	authCode := a.handler.CheckAppToken(req.Token, req.NameSpace, req.OperationUrl)

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
