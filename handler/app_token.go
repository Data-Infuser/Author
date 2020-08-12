package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"gitlab.com/promptech1/infuser-author/app/ctx"
	"gitlab.com/promptech1/infuser-author/constant"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/model"
)

type AppTokenHandler struct {
	Ctx *ctx.Context
}

func NewAppTokenHandler(ctx *ctx.Context) *AppTokenHandler {
	return &AppTokenHandler{
		Ctx: ctx,
	}
}

func (h *AppTokenHandler) CheckAppToken(tokenVal string, nameSpace string, operationUrl string) grpc_author.ApiAuthRes_Code {
	h.Ctx.Logger.Debug("CheckAppToken: ", tokenVal, nameSpace)
	var appId uint
	var tokenId uint
	var maxTraffic uint

	var operation model.Operation
	operation = model.Operation{
		EndPoint: operationUrl,
	}

	h.Ctx.Logger.Debug(operationUrl)

	h.Ctx.Orm.Table(operation.TableName()).
		Join(
			"INNER", operation.App.TableName(),
			"operation.app_id = app.id",
		).
		Where("app.name_space = ?", nameSpace).
		Get(&operation)

	h.Ctx.Logger.WithField(
		"operation", fmt.Sprintf("%+v", operation),
	).Debug("operation join result")

	app := model.App{NameSpace: nameSpace}

	nsKey := constant.REDIS_NS_KEY + nameSpace
	appInfo, err := h.Ctx.RedisDB.Get(nsKey, "string")
	if err != nil && err == redis.Nil {
		err = app.FindByNameSpace(h.Ctx.Orm)
		if err != nil {
			return grpc_author.ApiAuthRes_UNREGISTERED_SERVICE
		}
		// TODO : MaxTraffic 설정 처리 필요
		h.Ctx.RedisDB.Set(nsKey, fmt.Sprintf("%d:%d", app.Id, 100))

		appId = app.Id
		maxTraffic = 100 // TODO : MaxTraffic 설정 처리 필요
	} else {
		h.Ctx.Logger.Debug("find app in redis: ", appInfo)
		appInfoStr := appInfo
		splits := strings.Split(appInfoStr.(string), ":")

		appIDStr := splits[0]
		temp, _ := strconv.Atoi(appIDStr)
		appId = uint(temp)

		temp, _ = strconv.Atoi(splits[1])
		maxTraffic = uint(temp)
	}

	token := model.Token{Token: tokenVal}
	tokenKey := constant.REDIS_T_KEY + tokenVal
	tokenInfo, err := h.Ctx.RedisDB.Get(tokenKey, "uint")
	if err != nil && err == redis.Nil {
		if err = token.FindByToken(h.Ctx.Orm); err != nil {
			return grpc_author.ApiAuthRes_UNAUTHORIZED
		}
		h.Ctx.RedisDB.Set(tokenKey, token.Id)
		tokenId = token.Id
	} else {
		h.Ctx.Logger.Debug("find token in redis: ", tokenInfo)
		tokenId = tokenInfo.(uint)
	}

	appToken := model.AppToken{AppId: appId, TokenId: tokenId}
	err = appToken.FindOne(h.Ctx.Orm)
	if err != nil {
		return grpc_author.ApiAuthRes_UNAUTHORIZED
	}

	h.Ctx.Logger.Debug(fmt.Sprintf("AppToken ID: %d (appId: %d, tokenId: %d)", appToken.Id, appId, tokenId))
	var count uint
	trafficKey := fmt.Sprintf("traffic:%d", appToken.Id)
	countInfo, err := h.Ctx.RedisDB.Get(trafficKey, "uint")
	if err != nil && err == redis.Nil {
		count = uint(0)
		h.Ctx.RedisDB.SAdd(constant.REDIS_TRAFFIC_SET, trafficKey)
	} else {
		count = countInfo.(uint)
	}

	if maxTraffic >= count {
		// 인증키 활용횟수 Increment. 통계 저장은 별도 처리
		h.Ctx.RedisDB.Incr(trafficKey)
		return grpc_author.ApiAuthRes_VALID
	}

	return grpc_author.ApiAuthRes_LIMIT_EXCEEDED
}
