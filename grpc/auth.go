package server

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"gitlab.com/promptech1/infuser-author/constant"
	"gitlab.com/promptech1/infuser-author/handler"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"gitlab.com/promptech1/infuser-author/model"
)

type authServer struct {
	handler *handler.AuthHandler
}

func newAuthServer(handler *handler.AuthHandler) grpc_author.AuthServiceServer {
	return &authServer{handler: handler}
}

func (a authServer) Login(ctx context.Context, req *grpc_author.LoginReq) (*grpc_author.AuthRes, error) {
	ut := model.UserTokenRel{User: model.User{LoginId: req.LoginId}}

	if err := ut.FindByUserLoginId(a.handler.Ctx.Orm); err != nil {
		a.handler.Ctx.Logger.Info(err.Error())
		return &grpc_author.AuthRes{Code: grpc_author.AuthResult_NOT_REGISTERED}, nil
	}

	if _, err := model.ComparePasswords(ut.User.Password, req.Password); err != nil {
		a.handler.Ctx.Logger.Debug(err.Error())
		return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INVALID_PASSWORD}, nil
	}

	if ut.Token.Id == 0 || ut.Token.JwtExpiredAt == nil || time.Now().After(*ut.Token.JwtExpiredAt) {
		ut.Token.UserId = ut.User.Id

		// JWT 만료 시간 설정
		jwtExp := time.Now().Add(constant.JwtExpInterval)

		claims := &model.TokenClaims{
			LoginId: ut.User.LoginId, Email: ut.User.Email, Username: ut.User.Name,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: jwtExp.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		jwt, err := token.SignedString([]byte(constant.JwtSecret))
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			a.handler.Ctx.Logger.Info(err.Error())
			return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INVALID_PASSWORD}, nil
		}
		a.handler.Ctx.Logger.Debug(jwt)
		ut.Token.Jwt = jwt
		ut.Token.JwtExpiredAt = &jwtExp

		if len(ut.Token.RefreshToken) == 0 || time.Now().After(*ut.Token.RefreshTokenExpiredAt) {
			for {
				b := make([]byte, 32)
				rand.Read(b)
				refreshToken := fmt.Sprintf("%x", b)
				a.handler.Ctx.Logger.Debug(refreshToken)

				tokenInfo, err := model.FindUserTokenByRefreshToken(a.handler.Ctx.Orm, refreshToken)
				if err != nil {
					a.handler.Ctx.Logger.Debug(err)
				}
				if tokenInfo.Id == 0 {
					ut.Token.RefreshToken = refreshToken
					refreshExp := time.Now().Add(constant.RefreshTokenExpInterval)
					ut.Token.RefreshTokenExpiredAt = &refreshExp
					break
				}
			}
		}

		if err := ut.Token.Save(a.handler.Ctx.Orm); err != nil {
			a.handler.Ctx.Logger.Info(err.Error())
		}

		a.handler.Ctx.Logger.WithFields(logrus.Fields{
			"token detail": fmt.Sprintf("%+v", ut.Token),
		}).Debug("Token Info")
	}

	return &grpc_author.AuthRes{Code: grpc_author.AuthResult_VALID}, nil
}

func (a authServer) Refresh(ctx context.Context, req *grpc_author.RefreshTokenReq) (*grpc_author.AuthRes, error) {
	panic("implement me")
}
