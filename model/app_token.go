package model

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/promptech1/infuser-author/constant"
	errors "gitlab.com/promptech1/infuser-author/error"
	"xorm.io/xorm"
)

// App - Token relationship 정의.
type AppToken struct {
	Id        uint      `xorm:"pk autoincr"`
	AppId     uint      `xorm:"index"`
	TokenId   uint      `xorm:"index"`
	CreatedAt time.Time `xorm:"created"`

	App   App   `xorm:"- extends"`
	Token Token `xorm:"- extends"`
}

func (at *AppToken) FindOne(orm *xorm.Engine) error {
	found, err := orm.Get(at)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}

	if !found {
		return errors.NewWithCode(http.StatusNotFound, "token not found")
	}

	return nil
}

func (at *AppToken) KeyName() string {
	return fmt.Sprintf("%s%d:%d", constant.KeyAuth, at.TokenId, at.AppId)
}

func (at *AppToken) FindByAppAndToken(orm *xorm.Engine) error {
	found, err := orm.Get(at)

	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}

	if !found {
		return errors.NewWithCode(http.StatusNotFound, "appToken not found")
	}

	return nil
}
