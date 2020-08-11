package model

import "time"

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
