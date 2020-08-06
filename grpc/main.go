package grpc

import (
	"context"
	"github.com/golang/glog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"

	"gitlab.com/promptech1/infuser-author/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Run(ctx context.Context, network, address string) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	defer func() {
		if err := l.Close(); err != nil {
			glog.Errorf("Failed to clase %s %s: %v", network, address, err)
		}
	}()


	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	appTokenService := ctx.Value("appTokenService").(service.AppTokenService)

	// Token 기반의 인증 처리
	grpc_author.RegisterApiAuthServiceServer(s, newApiAuthServer(appTokenService))

	// App 서비스 등록/수정/삭제 처리
	grpc_author.RegisterAppTokenManagerServer(s, newAppTokenServer(appTokenService))

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	log.Printf("start gRPC grpc at %s", address)
	return s.Serve(l)
}