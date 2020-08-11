package model

import (
	"net/http"
	"time"

	errors "gitlab.com/promptech1/infuser-author/error"
	"xorm.io/xorm"
)

// App Api 서비스 관리 모델
type App struct {
	Id         uint      `xorm:"pk"`
	NameSpace  string    `xorm:"unique"`
	IsDel      bool      `xorm:index default 0`
	CreatedAt  time.Time `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  *time.Time
	Operations []Operation `xorm:"- extends"`
	Traffics   []Traffic   `xorm:"- extends"`
}

func (a *App) FindByNameSpace(orm *xorm.Engine) error {
	found, err := orm.Get(a)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}

	if !found {
		return errors.NewWithCode(http.StatusNotFound, "app not found")
	}

	return nil
}
