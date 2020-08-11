package main

import (
	"context"
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

	ctx := context.Background()

	var (
		err error
		a   *app.Application
	)

	a, err = app.New(ctx)
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
		_, err := ctx.RedisDB.LPop(constant.REDIS_TRAFFIC_QUEUE)
		if err != nil && err == redis.Nil {
			fmt.Println("Ignore stat ======= ")
			return
		}

		members, err := ctx.RedisDB.SMembers(constant.REDIS_TRAFFIC_SET)
		if err != nil && err == redis.Nil {
			fmt.Println("no members ======= ")
			return
		}

		var histories []model.AppTokenHistory

		for _, appTokenKey := range members {
			cntInfo, err := ctx.RedisDB.Get(appTokenKey, "uint")
			if err == nil {
				count := cntInfo.(uint)
				if count > 0 {
					ctx.RedisDB.Delete(appTokenKey)
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

		if len(histories) > 0 {
			ctx.Orm.Insert(&histories)
		}

		ctx.RedisDB.LPush(constant.REDIS_TRAFFIC_QUEUE, "1")

		time.Sleep(10 * 1000 * time.Millisecond)
		fmt.Println("Run Every min: ", time.Now().String())
	})
	c.Start()
}
