package service

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"gitlab.com/promptech1/infuser-author/database"
	"gitlab.com/promptech1/infuser-author/enum"
	"gitlab.com/promptech1/infuser-author/model"
	repo "gitlab.com/promptech1/infuser-author/repository"
	"strconv"
	"strings"
)

type AppTokenService interface {
	CheckAppToken(token string, nameSpace string) enum.AuthCode
	Regist(token string, nameSpace string)
}

type appTokenService struct {
	redisDB *database.RedisDB
	repo repo.AppTokenRepository
	appRepo repo.AppRepository
	tokenRepo repo.TokenRepository
}

func NewAppTokenService(
	repo repo.AppTokenRepository,
	appRepo repo.AppRepository,
	tokenRepo repo.TokenRepository,
	redisDB *database.RedisDB) AppTokenService {
	return &appTokenService{
		repo: repo,
		appRepo: appRepo,
		tokenRepo: tokenRepo,
		redisDB: redisDB,
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

func (s appTokenService) CheckAppToken(token string, nameSpace string) enum.AuthCode{
	var app *model.App
	var t *model.Token
	var appID uint
	var tokenID uint
	var maxTraffic uint

	nsKey := "ns:" + nameSpace
	appInfo, err := s.redisDB.Get(nsKey)
	if err != nil && err == redis.Nil {
		if app, err = s.appRepo.FindByNameSpace(nameSpace); err != nil {
			return enum.UNREGISTERED_SERVICE
		}

		glog.Info("find app in db ============")
		s.redisDB.Set(nsKey, fmt.Sprintf("%d:%d", app.ID, app.MaxTraffic))
		appID = app.ID
		maxTraffic = app.MaxTraffic
	} else {
		glog.Info("find app in redis ============")
		splits := strings.Split(appInfo, ":")
		appIDStr := splits[0]
		temp, _ := strconv.Atoi(appIDStr)
		appID = uint(temp)

		temp, _ = strconv.Atoi(splits[1])
		maxTraffic = uint(temp)
	}

	tKey := "t:" + token
	tokenIDStr, err := s.redisDB.Get(tKey)
	if err != nil && err == redis.Nil {
		if t, err = s.tokenRepo.FindByToken(token); err != nil {
			return enum.UNAUTHORIZED
		}
		s.redisDB.Set(tKey, t.ID)
		tokenID = t.ID
	} else {
		glog.Info("find token in redis ============")
		temp, _ := strconv.Atoi(tokenIDStr)
		tokenID = uint(temp)
	}

	appToken := s.repo.Find(appID, tokenID)
	if appToken != nil {
		count := s.repo.FindTodayUsage(appToken)
		if maxTraffic >= count {
			// TODO: API 인증 횟수 추가 처리 필요함
			return enum.VALID
		}

		return enum.LIMIT_EXCEEDED
	}

	return enum.UNAUTHORIZED
}
