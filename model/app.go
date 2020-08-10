package model

import "time"

// App Api 서비스 관리 모델
type App struct {
	Id uint `xorm:"pk"`
	NameSpace string `xorm:"unique"`
	IsDel bool `xorm:default 0`

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time
}
