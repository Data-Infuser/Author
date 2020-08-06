package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"gitlab.com/promptech1/infuser-author/database"
	"gitlab.com/promptech1/infuser-author/grpc"
	repo "gitlab.com/promptech1/infuser-author/repository"
	"gitlab.com/promptech1/infuser-author/service"
)

var (
	addr    = flag.String("addr", ":9090", "endpoint of the gRPC grpc")
	network    = flag.String("network", "tcp", `one of "tcp" or "unix". Must be consistent to -endpoint`)
)

func main() {
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()

	db := database.ConnDB()
	defer db.Close()
	ctx = context.WithValue(ctx, "db", db)

	redisDB := database.ConnRedis(ctx)
	ctx = context.WithValue(ctx, "redisDB", redisDB)

	tokenRepo := repo.NewTokenRepository(db)
	appRepo := repo.NewAppRepository(db)
	appTokenRepo := repo.NewAppTokenRepository(db)
	//userRepo := repo.NewUserRepository(db)

	ctx = context.WithValue(ctx, "tokenRepo", tokenRepo)
	ctx = context.WithValue(ctx, "appRepo", appRepo)
	ctx = context.WithValue(ctx, "appTokenRepo", appTokenRepo)

	//userService := service.NewUserService(userRepo)
	appTokenService := service.NewAppTokenService(ctx)

	ctx = context.WithValue(ctx, "appTokenService", appTokenService)

	if err := grpc.Run(ctx, *network, *addr); err != nil {
		glog.Fatal(err)
	}
}