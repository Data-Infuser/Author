package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/promptech1/infuser-author/constant"
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
		if _, err := orm.Update(ut); err != nil {
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

func CheckRefreshToken(orm *xorm.Engine, refreshToken string) (bool, error) {
	ut := &UserToken{RefreshToken: refreshToken}
	return orm.Get(ut)
}

func (ut *UserToken) FindUserTokenByRefreshToken(orm *xorm.Engine) error {
	if _, err := orm.Get(ut); err != nil {
		return err
	}

	return nil
}
