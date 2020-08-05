package repo

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/promptech1/infuser-author/model"
)

type TokenRepository interface {
	Create(token string) *model.Token
	FindByToken(token string) (*model.Token, error)
}

type tokenRepositoryDB struct {
	DB *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepositoryDB{DB: db}
}

func (r *tokenRepositoryDB) Create(token string) *model.Token {
	t := &model.Token{Token: token}
	r.DB.Create(t)

	return t
}

func (r *tokenRepositoryDB) FindByToken(token string) (*model.Token, error) {
	t := model.Token{}

	if err := r.DB.Where("token = ?", token).First(&t).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return &t, nil
}