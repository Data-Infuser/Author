package model

import "time"

type Role struct {
	Id        uint `xorm:"pk autoincr"`
	Name      string
	CreatedAt time.Time `xorm:"created"`
}
