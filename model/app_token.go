package model

import "github.com/jinzhu/gorm"

// App - Token relationship 정의. 현재 모델에 대한 Migration 수행 이후 각 모델의 many-to-many 관계 정의 필요
type AppToken struct {
	gorm.Model
	AppID uint
	App App
	TokenID uint
	Token Token

	Histories []AppTokenHistory
}
