package model

import "github.com/jinzhu/gorm"

// App Api 서비스 관리 모델
type App struct {
	gorm.Model
	NameSpace string `gorm:"size:191;unique_index"`
	MaxTraffic uint // 1일 허용 트래픽 횟수 (향후 트래픽 정책에 따라 1:N 등의 별도 분리 필요)

	// Many-to-Many 관계는 AppToken 모델을 먼저 정의 후 기입함(DB Migration 순서 확인필요)
	Tokens         []Token `gorm:"many2many:app_tokens;"`
}
