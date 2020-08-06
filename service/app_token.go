package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"gitlab.com/promptech1/infuser-author/database"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/model"
	repo "gitlab.com/promptech1/infuser-author/repository"
	"strconv"
	"strings"
)

const TrafficSet = "usage-set"

type AppTokenService interface {
	CheckAppToken(token string, nameSpace string) grpc_author.ApiAuthRes_Code
	Regist(token string, nameSpace string)
}

type appTokenService struct {
	redisDB *database.RedisDB
	repo repo.AppTokenRepository
	appRepo repo.AppRepository
	tokenRepo repo.TokenRepository
}

func NewAppTokenService(ctx context.Context) AppTokenService {
	return &appTokenService{
		repo: ctx.Value("appTokenRepo").(repo.AppTokenRepository),
		appRepo: ctx.Value("appRepo").(repo.AppRepository),
		tokenRepo: ctx.Value("tokenRepo").(repo.TokenRepository),
		redisDB: ctx.Value("redisDB").(*database.RedisDB),
	}
}

func (s appTokenService) Regist(token string, nameSpace string) {
	app, _ := s.appRepo.FindByNameSpace(nameSpace)
	t, _ := s.tokenRepo.FindByToken(token)
	if t == nil {
		t = s.tokenRepo.Create(token)
	}

	s.repo.Create(app.ID, t.ID)
}

func (s appTokenService) CheckAppToken(token string, nameSpace string) grpc_author.ApiAuthRes_Code{
	glog.Info("CheckAppToken: ", token, nameSpace)
	var app *model.App
	var t *model.Token
	var appID uint
	var tokenID uint
	var maxTraffic uint

	// APP 정보 조회 (1.redis, 2.RDB)
	nsKey := "ns:" + nameSpace
	if appInfo, err := s.redisDB.Get(nsKey, "string"); err != nil && err == redis.Nil {
		glog.Info("find app in db ============")
		if app, err = s.appRepo.FindByNameSpace(nameSpace); err != nil {
			return grpc_author.ApiAuthRes_UNREGISTERED_SERVICE
		}

		// Redis에 저장하는 APP 정보는 ID와 최대 트래픽(delimiter를 ':' 로 사용)
		s.redisDB.Set(nsKey, fmt.Sprintf("%d:%d", app.ID, app.MaxTraffic))
		appID = app.ID
		maxTraffic = app.MaxTraffic
	} else {
		glog.Infof("find app in redis: %s", appInfo)
		appInfoStr := appInfo.(string)
		splits := strings.Split(appInfoStr, ":")

		appIDStr := splits[0]
		temp, _ := strconv.Atoi(appIDStr)
		appID = uint(temp)

		temp, _ = strconv.Atoi(splits[1])
		maxTraffic = uint(temp)
	}

	// Token 정보 조회(1.redis, 2.RDB)
	tKey := "t:" + token
	if tokenIDStr, err := s.redisDB.Get(tKey, "uint"); err != nil && err == redis.Nil {
		if t, err = s.tokenRepo.FindByToken(token); err != nil {
			return grpc_author.ApiAuthRes_UNAUTHORIZED
		}
		s.redisDB.Set(tKey, t.ID)
		tokenID = t.ID
	} else {
		glog.Info("find token in redis: ", tokenIDStr)
		tokenID = tokenIDStr.(uint)
	}

	// App-Token 정보 조회
	appToken := s.repo.Find(appID, tokenID)
	if appToken != nil {
		glog.Infof("AppToken ID: %d (appId: %d, tokenId: %d)", appToken.ID, appID, tokenID)
		trafficKey := fmt.Sprintf("traffic:%d", appToken.ID)

		count, err := s.redisDB.Get(trafficKey, "uint")
		if err != nil && err == redis.Nil {
			count = uint(0)
			s.redisDB.SAdd(TrafficSet, trafficKey)
		}

		//count := s.repo.FindTodayUsage(appToken)
		if maxTraffic >= count.(uint) {
			// 인증키 활용횟수 Increment. 통계 저장은 별도 처리
			s.redisDB.Incr(trafficKey)
			return grpc_author.ApiAuthRes_VALID
		}

		return grpc_author.ApiAuthRes_LIMIT_EXCEEDED
	}

	return grpc_author.ApiAuthRes_UNAUTHORIZED
}
