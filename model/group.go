package model

import "time"

type Group struct {
	Id        uint `xorm:"pk autoincr"`
	Name      string
	CreatedAt time.Time  `xorm:"created"`
	UpdatedAt time.Time  `xorm:"updated"`
	DeletedAt *time.Time `xorm:"deleted"`
}
