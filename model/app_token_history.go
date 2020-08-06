package model

import "time"

type AppTokenHistory struct {
	ID        uint `gorm:"primary_key"`
	AppTokenID uint `gorm:"index"`
	AppToken AppToken
	CallTraffic uint
	CreatedAt time.Time `gorm:"index"`
}