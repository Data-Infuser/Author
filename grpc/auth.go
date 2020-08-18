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
	"gitlab.com/promptech1/infuser-author/model/relations"
)

type authServer struct {
	handler *handler.AuthHandler
}

func newAuthServer(handler *handler.AuthHandler) grpc_author.AuthServiceServer {
	return &authServer{handler: handler}
}

func (a *authServer) Login(ctx context.Context, req *grpc_author.LoginReq) (*grpc_author.AuthRes, error) {
	utr := relations.UserTokenRel{User: model.User{LoginId: req.LoginId}}

	// 회원 조회
	if err := utr.FindByUserLoginId(a.handler.Ctx.Orm); err != nil {
		a.handler.Ctx.Logger.Info(err.Error())
		return &grpc_author.AuthRes{Code: grpc_author.AuthResult_NOT_REGISTERED}, nil
	}

	// 비밀번호 확인
	if _, err := model.ComparePasswords(utr.User.Password, req.Password); err != nil {
		a.handler.Ctx.Logger.Debug(err.Error())
		return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INVALID_PASSWORD}, nil
	}

	a.handler.Ctx.Logger.WithFields(logrus.Fields{
		"UserTokenRel": fmt.Sprintf("%+v", utr),
	}).Debug("Token Info")

	if utr.Token.Id == 0 || utr.Token.JwtExpiredAt == nil || time.Now().After(*utr.Token.JwtExpiredAt) {
		utr.Token.UserId = utr.User.Id

		// JWT 만료 시간 설정
		jwtExp := time.Now().Add(constant.JwtExpInterval)

		claims := &model.TokenClaims{
			LoginId: utr.User.LoginId, Email: utr.User.Email, Username: utr.User.Name,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: jwtExp.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		jwt, err := token.SignedString([]byte(constant.JwtSecret))
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			a.handler.Ctx.Logger.Info(err.Error())
			return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INTERNAL_EXCEPTION}, nil
		}
		a.handler.Ctx.Logger.Debug(jwt)
		utr.Token.Jwt = jwt
		utr.Token.JwtExpiredAt = &jwtExp

		if len(utr.Token.RefreshToken) == 0 || time.Now().After(*utr.Token.RefreshTokenExpiredAt) {
			refreshToken, err := a.genRefreshToken()
			if err != nil {
				a.handler.Ctx.Logger.Info(err.Error())
				return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INTERNAL_EXCEPTION}, nil
			}

			utr.Token.SetRefreshToken(refreshToken)
		}

		if err := utr.Token.Save(a.handler.Ctx.Orm); err != nil {
			a.handler.Ctx.Logger.Info(err.Error())
			return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INTERNAL_EXCEPTION}, nil
		}

		a.handler.Ctx.Logger.WithFields(logrus.Fields{
			"token detail": fmt.Sprintf("%+v", utr.Token),
		}).Debug("Token Info")
	}

	return utr.Token.GetValidGrpcRes()
}

func (a *authServer) Auth(ctx context.Context, req *grpc_author.JwtReq) (*grpc_author.AuthRes, error) {
	ut := model.UserToken{Jwt: req.Jwt}
	err := ut.FindUserToken(a.handler.Ctx.Orm)

	if err != nil {
		a.handler.Ctx.Logger.Info(err.Error())
		return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INTERNAL_EXCEPTION}, nil
	}

	if ut.Id == 0 {
		return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INVALID_TOKEN}, nil
	}

	return ut.GetValidGrpcRes()
}

func (a *authServer) Refresh(ctx context.Context, req *grpc_author.RefreshTokenReq) (*grpc_author.AuthRes, error) {
	ut := model.UserToken{RefreshToken: req.RefreshToken}
	err := ut.FindUserToken(a.handler.Ctx.Orm)

	if err != nil {
		a.handler.Ctx.Logger.Info(err.Error())
		return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INTERNAL_EXCEPTION}, nil
	}

	if ut.Id == 0 {
		return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INVALID_TOKEN}, nil
	}

	refreshToken, err := a.genRefreshToken()
	if err != nil {
		a.handler.Ctx.Logger.Info(err.Error())
		return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INTERNAL_EXCEPTION}, nil
	}

	ut.SetRefreshToken(refreshToken)
	if err := ut.Save(a.handler.Ctx.Orm); err != nil {
		a.handler.Ctx.Logger.Info(err.Error())
		return &grpc_author.AuthRes{Code: grpc_author.AuthResult_INTERNAL_EXCEPTION}, nil
	}

	return ut.GetValidGrpcRes()
}

func (a *authServer) genRefreshToken() (string, error) {
	for {
		b := make([]byte, 32)
		rand.Read(b)
		refreshToken := fmt.Sprintf("%x", b)
		a.handler.Ctx.Logger.Debug(refreshToken)

		has, err := model.CheckRefreshToken(a.handler.Ctx.Orm, refreshToken)
		if err != nil {
			a.handler.Ctx.Logger.Debug(err)
			return "", err
		}
		if !has {
			return refreshToken, nil
		}
	}
}
