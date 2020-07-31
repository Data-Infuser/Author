package grpc

import (
	"context"
	"github.com/golang/glog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"gitlab.com/promptech1/infuser-author/database"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"

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

	userRepo := repo.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	grpc_author.RegisterTokenManagerServer(s, newTokenServer(tokenService))
	grpc_author.RegisterUserManagerServer(s, newUserServer(userService))

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	log.Printf("start gRPC grpc at %s", address)
	return s.Serve(l)
}