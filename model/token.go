package model

import (
	"net/http"
	"time"

	"gitlab.com/promptech1/infuser-author/constant"
	errors "gitlab.com/promptech1/infuser-author/error"
	"xorm.io/xorm"
)

// Token : API 인증 토큰 관리 모델
type Token struct {
	Id    uint   `xorm:"pk autoincr"`
	Token string `xorm:"unique"`
	IsDel bool   `xorm:index default 0`

	CreatedAt time.Time `xorm:"created"`
	DeletedAt *time.Time
}

func (t *Token) KeyName() string {
	return constant.KeyToken + t.Token
}

func (t *Token) FindByToken(orm *xorm.Engine) error {
	found, err := orm.Get(t)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}

	if !found {
		return errors.NewWithCode(http.StatusNotFound, "token not found")
	}

	return nil
}
