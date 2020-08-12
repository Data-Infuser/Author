package model

import (
	"net/http"
	"time"

	"gitlab.com/promptech1/infuser-author/constant"
	errors "gitlab.com/promptech1/infuser-author/error"
	"xorm.io/xorm"
)

type Operation struct {
	Id        uint `xorm:"pk"`
	AppId     uint `xorm:"index"`
	EndPoint  string
	IsDel     bool      `xorm:index default 0`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt *time.Time

	App App `xorm:"- extends"`
}

func (Operation) TableName() string {
	return "operation"
}

func (o *Operation) KeyName() string {
	return constant.KeyOperation + o.EndPoint
}

func (o *Operation) FindByEndPoint(orm *xorm.Engine) error {
	found, err := orm.Get(o)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}

	if !found {
		return errors.NewWithCode(http.StatusNotFound, "operation not found")
	}

	return nil
}
