package model

import (
	"net/http"
	"time"

	"gitlab.com/promptech1/infuser-author/constant"
	"gitlab.com/promptech1/infuser-author/database"
	errors "gitlab.com/promptech1/infuser-author/error"
	"xorm.io/xorm"
)

type Operation struct {
	Id        uint `xorm:"pk"`
	AppId     uint `xorm:"index"`
	EndPoint  string
	IsDel     bool       `xorm:index default 0`
	Version   int        `xorm:"version"`
	CreatedAt time.Time  `xorm:"created"`
	UpdatedAt time.Time  `xorm:"updated"`
	DeletedAt *time.Time `xorm:"deleted index"`

	App App `xorm:"- extends"`
}

func (Operation) TableName() string {
	return "operation"
}

func (o *Operation) KeyName() string {
	return constant.KeyOperation + o.EndPoint
}

func (o *Operation) FindOperation(orm *xorm.Engine) error {
	found, err := orm.Get(o)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}

	if !found {
		return errors.NewWithCode(http.StatusNotFound, "operation not found")
	}

	return nil
}

func (o *Operation) SetRedis(rdb *database.RedisDB) {
	rdb.Set(o.KeyName(), o.Id)
}

func (o *Operation) DelRedis(rdb *database.RedisDB) {
	rdb.Delete(o.KeyName())
}

func (o *Operation) Update(orm *xorm.Engine) error {
	if _, err := orm.ID(o.Id).Update(o); err != nil {
		return err
	}

	return nil
}

func (o *Operation) Delete(orm *xorm.Engine) error {
	sql := "UPDATE operation SET deleted_at = ? WHERE id = ?"
	if _, err := orm.Exec(sql, time.Now(), o.Id); err != nil {
		return err
	}

	return nil
}

func FindOperationsByApp(orm *xorm.Engine, appId uint) ([]Operation, error) {
	operations := []Operation{}

	err := orm.Where("app_id = ?", appId).Find(&operations)

	if err != nil {
		return nil, errors.New("database error; " + err.Error())
	}

	return operations, nil
}
