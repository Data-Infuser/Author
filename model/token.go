package model

import "github.com/jinzhu/gorm"

type Token struct {
	gorm.Model
	Token string `gorm:"unique_index"`
}