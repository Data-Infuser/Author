package grpc

import (
	"context"
	"gitlab.com/promptech1/infuser-author/gen/gitlab.com/promptech1/infuser-author/gen"
	"gitlab.com/promptech1/infuser-author/service"
)

type tokenServer struct {
	tokenService service.TokenService
}

func newTokenServer(tokenService service.TokenService) gen.TokenManagerServer {
	return &tokenServer{
		tokenService: tokenService,
	}
}

func (s *tokenServer) Search(ctx context.Context, req *gen.TokenReq) (*gen.TokenRes, error) {
	token := s.tokenService.CheckToken(req.Token)

	res := &gen.TokenRes{
		Token: token.Token,
		Status: gen.TokenRes_VALID,
	}

	return res, nil
}