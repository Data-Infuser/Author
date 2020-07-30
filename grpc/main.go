package grpc

import (
	"context"
	"github.com/golang/glog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"gitlab.com/promptech1/infuser-author/database"
	"gitlab.com/promptech1/infuser-author/gen/gitlab.com/promptech1/infuser-author/gen"
	repo "gitlab.com/promptech1/infuser-author/repository"
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

	db := database.ConnDB()
	defer db.Close()

	tokenRepo := repo.NewTokenRepository(db)
	tokenService := service.NewTokenService(tokenRepo)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	gen.RegisterTokenManagerServer(s, newTokenServer(tokenService))

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	log.Printf("start gRPC grpc at %s", address)
	return s.Serve(l)
}