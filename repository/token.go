package repo

import (
	"gitlab.com/promptech1/infuser-author/model"
	"xorm.io/xorm"
)

type TokenRepository interface {
	Create(token string) *model.Token
	FindByToken(token string) (*model.Token, error)
}

type tokenRepositoryDB struct {
	DB *xorm.Engine
}

func NewTokenRepository(db *xorm.Engine) TokenRepository {
	return &tokenRepositoryDB{DB: db}
}

func (r *tokenRepositoryDB) Create(token string) *model.Token {
	t := &model.Token{Token: token}
	r.DB.Insert(t)

	return t
}

func (r *tokenRepositoryDB) FindByToken(token string) (*model.Token, error) {
	t := model.Token{Token: token}

	if _, err := r.DB.Get(&t); err != nil {
		return nil, err
	}

	return &t, nil
}