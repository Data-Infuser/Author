package grpc

import (
	"context"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/service"
)

type apiAuthServer struct {
	appTokenService service.AppTokenService
}

func (a *apiAuthServer) Auth(ctx context.Context, req *grpc_author.ApiAuthReq) (*grpc_author.ApiAuthRes, error) {
	authCode := a.appTokenService.CheckAppToken(req.Token, req.NameSpace)

	res := &grpc_author.ApiAuthRes{
		Code: authCode,
	}

	return res, nil
}

func newApiAuthServer(appTokenService service.AppTokenService) grpc_author.ApiAuthServiceServer {
	return &apiAuthServer{
		appTokenService: appTokenService,
	}
}