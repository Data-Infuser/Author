package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"gitlab.com/promptech1/infuser-author/constant"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"xorm.io/xorm"
)

// UseToken : 회원 인증 토큰 관리 모델
type UserToken struct {
	Id                    uint   `xorm:"pk autoincr"`
	UserId                uint   `xorm:"index"`
	Jwt                   string `xorm:"unique"`
	RefreshToken          string `xorm:"unique"`
	JwtExpiredAt          *time.Time
	RefreshTokenExpiredAt *time.Time
	CreatedAt             time.Time `xorm:"created"`
	UpdatedAt             time.Time `xorm:"updated"`
}

type TokenClaims struct {
	Id       uint   `json:"id"`
	LoginId  string `json:"loginId"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func (ut *UserToken) Save(orm *xorm.Engine) error {
	if ut.Id == 0 {
		if _, err := orm.Insert(ut); err != nil {
			return err
		}
	} else {
		if _, err := orm.ID(ut.Id).Update(ut); err != nil {
			return err
		}
	}

	return nil
}

func (ut *UserToken) SetRefreshToken(refreshToken string) {
	ut.RefreshToken = refreshToken
	refreshExp := time.Now().Add(constant.RefreshTokenExpInterval)
	ut.RefreshTokenExpiredAt = &refreshExp
}

func (ut *UserToken) GetValidGrpcRes() (*grpc_author.AuthRes, error) {
	var expiresIn *timestamp.Timestamp
	var refreshTokenExpiresIn *timestamp.Timestamp
	var err error

	if ut.JwtExpiredAt != nil {
		if expiresIn, err = ptypes.TimestampProto(*ut.JwtExpiredAt); err != nil {
			return nil, err
		}
	}

	if ut.RefreshTokenExpiredAt != nil {
		if refreshTokenExpiresIn, err = ptypes.TimestampProto(*ut.RefreshTokenExpiredAt); err != nil {
			return nil, err
		}
	}

	// TODO: BloopRPC에서 timestamp 값이 포함된 경우 response의 출력오류 발생(실제 값은 정상 입력되어있음)
	return &grpc_author.AuthRes{
		Code:                  grpc_author.AuthResult_VALID,
		Jwt:                   ut.Jwt,
		RefreshToken:          ut.RefreshToken,
		ExpiresIn:             expiresIn,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}, nil
}

func (ut *UserToken) FindUserToken(orm *xorm.Engine) error {
	if _, err := orm.Get(ut); err != nil {
		return err
	}

	return nil
}

func CheckRefreshToken(orm *xorm.Engine, refreshToken string) (bool, error) {
	ut := &UserToken{RefreshToken: refreshToken}
	return orm.Get(ut)
}
