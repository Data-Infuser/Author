package model

import "github.com/jinzhu/gorm"

// Token : API 인증 토큰 관리 모델
type Token struct {
	gorm.Model
	Token string `gorm:"unique_index"`

	// Many-to-Many 관계는 AppToken 모델을 먼저 정의 후 기입함(DB Migration 순서 확인필요)
	Apps         []App `gorm:"many2many:app_tokens;"`
}
