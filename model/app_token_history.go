package model

import "time"

type AppTokenHistory struct {
	Id          uint `xorm:"pk autoincr"`
	AppTokenId  uint `xorm:"index"`
	OperationId uint `xorm:"index"`
	CallTraffic uint
	CreatedAt   time.Time `xorm:"created"`

	AppToken  AppToken  `xorm:"- extends"`
	Operation Operation `xorm:"- extends"`
}
