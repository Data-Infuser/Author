package repo

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/promptech1/infuser-author/model"
	"log"
)

type TokenRepository interface {
	FindByToken(tokenVal string) model.Token
}

type tokenRepositoryDB struct {
	DB *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepositoryDB{DB: db}
}

func (r *tokenRepositoryDB) FindByToken(tokenVal string) model.Token {
	t := model.Token{}

	if r.DB.Where("token = ?", tokenVal).First(&t).RecordNotFound() {
		log.Printf("cannot found token info for '%s'", tokenVal)
	}

	return t
}