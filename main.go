package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"github.com/robfig/cron/v3"
	"gitlab.com/promptech1/infuser-author/app"
	"gitlab.com/promptech1/infuser-author/app/ctx"
	"gitlab.com/promptech1/infuser-author/constant"
	"gitlab.com/promptech1/infuser-author/model"
)

var (
	addr    = flag.String("addr", ":9090", "endpoint of the gRPC grpc")
	network = flag.String("network", "tcp", `one of "tcp" or "unix". Must be consistent to -endpoint`)
)

func main() {
	flag.Parse()
	defer glog.Flush()

	var (
		err error
		a   *app.Application
	)

	a, err = app.New()
	if err != nil {
		log.Fatal("error ", os.Args[0]+" initialization error: "+err.Error())
		os.Exit(1)
	}

	go runCron(a.Ctx)

	a.Run(*network, *addr)

	glog.Info("Run Grpc Service =======")
}

func runCron(ctx *ctx.Context) {
	c := cron.New()
	c.AddFunc("* * * * *", func() {

		_, err := ctx.RedisDB.LPop(ctx.Context, constant.REDIS_TRAFFIC_QUEUE).Result()
		if err != nil && err == redis.Nil {
			fmt.Println("Ignore stat ======= ")
			return
		}

		members, err := ctx.RedisDB.SMembers(ctx.Context, constant.REDIS_TRAFFIC_SET).Result()
		if err != nil && err == redis.Nil {
			fmt.Println("no members ======= ")
			return
		}

		var histories []model.AppTokenHistory

		for _, appTokenKey := range members {
			cntInfo, err := ctx.RedisDB.Get(ctx.Context, appTokenKey).Result()
			if err == nil {
				temp, err := strconv.ParseUint(cntInfo, 10, 32)
				if err == nil {
					count := uint(temp)
					if count > 0 {
						ctx.RedisDB.Del(ctx.Context, appTokenKey).Result()
						temp1 := strings.Split(appTokenKey, ":")[1]
						temp, _ := strconv.Atoi(temp1)
						appTokenId := uint(temp)

						histories = append(histories, model.AppTokenHistory{
							AppTokenId:  appTokenId,
							CallTraffic: count,
						})
					}
				}
			}
		}

		if len(histories) > 0 {
			ctx.Orm.Insert(&histories)
		}

		ctx.RedisDB.LPush(ctx.Context, constant.REDIS_TRAFFIC_QUEUE, "1").Result()

		time.Sleep(10 * 1000 * time.Millisecond)
		fmt.Println("Run Every min: ", time.Now().String())
	})
	c.Start()
}
