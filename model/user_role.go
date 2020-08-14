package model

import "time"

type UserRole struct {
	Id        uint      `xorm:"pk autoincr"`
	UserId    uint      `xorm:"index"`
	RoleId    uint      `xorm:"index"`
	CreatedAt time.Time `xorm:"created"`

	User User `xorm:"- extends"`
	Role Role `xorm:"- extends"`
}
