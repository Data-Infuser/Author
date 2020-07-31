package grpc

import (
	"context"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/service"

	"log"
)

type tokenServer struct {
	tokenService service.TokenService
}

func newTokenServer(tokenService service.TokenService) grpc_author.TokenManagerServer {
	return &tokenServer{
		tokenService: tokenService,
	}
}

func (s *tokenServer) Search(ctx context.Context, req *grpc_author.TokenReq) (*grpc_author.TokenRes, error) {
	token := s.tokenService.CheckToken(req.Token)

	res := &grpc_author.TokenRes{
		Token: token.Token,
		Status: grpc_author.TokenRes_VALID,
	}

	log.Printf("Search result: %s", res.Token)
	log.Printf("Search result: %v", res.Status)

	return res, nil
}