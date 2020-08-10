package model

import "time"

// Token : API 인증 토큰 관리 모델
type Token struct {
	Id    uint   `xorm:"pk;autoincr"`
	Token string `xorm:"unique"`
	IsDel bool   `xorm:default 0`

	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time
}
