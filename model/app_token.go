package model

import "time"

// App - Token relationship 정의.
type AppToken struct {
	Id uint `xorm:"pk;autoincr"`
	AppId uint `xorm:"index"`
	TokenId uint `xorm:"index"`
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time

	App App `xorm:"- extends"`
	Token Token `xorm:"- extends"`
}
