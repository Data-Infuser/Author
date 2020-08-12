package server

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"gitlab.com/promptech1/infuser-author/app/ctx"
	"gitlab.com/promptech1/infuser-author/handler"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"google.golang.org/grpc"
)

// Server is an main application object that shared (read-only) to application modules
type Server struct {
	ctx        *ctx.Context
	context    context.Context
	grpcServer *grpc.Server
}

// New constructor
func New(c *ctx.Context, context context.Context) *Server {
	s := new(Server)
	s.ctx = c
	s.context = context
	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	return s
}

func (s *Server) Run(network, address string) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	appTokenHandler := handler.NewAppTokenHandler(s.ctx)
	appHandler := handler.NewAppHandler(s.ctx)

	// Token 기반의 인증 처리
	grpc_author.RegisterApiAuthServiceServer(s.grpcServer, newApiAuthServer(appTokenHandler))

	grpc_author.RegisterAppManagerServer(s.grpcServer, newAppManagerServer(appHandler))

	go func() {
		defer s.grpcServer.GracefulStop()
		<-s.context.Done()
	}()

	s.ctx.Logger.Info("start gRPC grpc at ", address)
	return s.grpcServer.Serve(l)

}
