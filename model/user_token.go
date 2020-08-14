package model

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	errors "gitlab.com/promptech1/infuser-author/error"
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

type UserTokenRel struct {
	User  User      `xorm:"extends"`
	Token UserToken `xorm:"extends"`
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

func (ut *UserTokenRel) FindByUserLoginId(orm *xorm.Engine) error {
	found, err := orm.Table("user").Join(
		"LEFT OUTER", "user_token",
		"user.id = user_token.user_id",
	).Where("user.login_id = ?", ut.User.LoginId).Get(ut)

	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}

	if !found {
		return errors.NewWithCode(http.StatusNotFound, "UserToken not found")
	}

	return nil
}

func FindUserTokenByRefreshToken(orm *xorm.Engine, refreshToken string) (*UserToken, error) {
	ut := UserToken{RefreshToken: refreshToken}
	if _, err := orm.Get(&ut); err != nil {
		return nil, err
	}

	return &ut, nil
}
