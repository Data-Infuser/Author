package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"github.com/robfig/cron/v3"
	"gitlab.com/promptech1/infuser-author/constant"
	"gitlab.com/promptech1/infuser-author/database"
	"gitlab.com/promptech1/infuser-author/grpc"
	"gitlab.com/promptech1/infuser-author/model"
	repo "gitlab.com/promptech1/infuser-author/repository"
	"gitlab.com/promptech1/infuser-author/service"
	"strconv"
	"strings"
	"time"
)

var (
	addr    = flag.String("addr", ":9090", "endpoint of the gRPC grpc")
	network    = flag.String("network", "tcp", `one of "tcp" or "unix". Must be consistent to -endpoint`)
)

const trafficQueueName = "traffic-history"

func main() {
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()

	db := database.ConnDB()
	defer db.Close()
	ctx = context.WithValue(ctx, constant.DB, db)

	redisDB := database.ConnRedis(ctx)
	ctx = context.WithValue(ctx, constant.REDIS_DB, redisDB)

	tokenRepo := repo.NewTokenRepository(db)
	appRepo := repo.NewAppRepository(db)
	appTokenRepo := repo.NewAppTokenRepository(db)
	appTokenHistoryRepo := repo.NewAppTokenHistoryRepository(db)
	//userRepo := repo.NewUserRepository(db)

	ctx = context.WithValue(ctx, "tokenRepo", tokenRepo)
	ctx = context.WithValue(ctx, "appRepo", appRepo)
	ctx = context.WithValue(ctx, "appTokenRepo", appTokenRepo)
	ctx = context.WithValue(ctx, "appTokenHistoryRepo", appTokenHistoryRepo)

	//userService := service.NewUserService(userRepo)
	appTokenService := service.NewAppTokenService(ctx)

	ctx = context.WithValue(ctx, "appTokenService", appTokenService)

	go runCron(ctx)

	if err := grpc.Run(ctx, *network, *addr); err != nil {
		glog.Fatal(err)
	}

	glog.Info("Run Grpc Service =======")
}

func runCron(ctx context.Context) {
	c := cron.New()
	c.AddFunc("* * * * *", func() {
		redisDB := ctx.Value("redisDB").(*database.RedisDB)
		if _, err := redisDB.LPop(trafficQueueName); err != nil && err == redis.Nil {
			fmt.Println("Ignore stat ======= ")
			return
		}

		members, err := redisDB.SMembers(service.TrafficSet)
		if err != nil && err == redis.Nil {
			fmt.Println("no members ======= ")
			return
		}

		var histories []model.AppTokenHistory

		for _, appTokenKey := range members {
			cnt, err := redisDB.Get(appTokenKey, "uint")
			if err == nil && cnt.(uint) > 0 {
				redisDB.Delete(appTokenKey)

				temp1 := strings.Split(appTokenKey, ":")[1]
				temp, _ := strconv.Atoi(temp1)
				appTokenID := uint(temp)

				histories = append(histories, model.AppTokenHistory{
					AppTokenID: appTokenID,
					CallTraffic: cnt.(uint),
					CreatedAt: time.Now(),
				})
			}
		}

		if len(histories) > 0 {
			appTokenHistoryRepo := ctx.Value("appTokenHistoryRepo").(repo.AppTokenHistoryRepository)
			appTokenHistoryRepo.Create(histories)
		}

		redisDB.LPush(trafficQueueName, "1")

		fmt.Println("Run Every min: ", time.Now().String())
	})
	c.Start()
}