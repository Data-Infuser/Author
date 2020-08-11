package model

import "time"

type Traffic struct {
	Id        uint   `xorm:"pk autoincr"`
	AppId     uint   `xorm:"index index(with_seq)"`
	Unit      string `xorm:"varchar(10) index default 'd'"` //트래픽 단위(min:분, hour:시간, day:1일, month:1달)
	Val       uint
	Seq       uint      `xorm:"index(with_seq)"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt *time.Time

	App App `xorm:"- extends"`
}
