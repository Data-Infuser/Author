package service

import (
	"gitlab.com/promptech1/infuser-author/model"
	"gitlab.com/promptech1/infuser-author/repository"
)

type TokenService interface {
	CheckToken(token string) model.Token
}

type tokenService struct {
	repo repo.TokenRepository
}

func NewTokenService(repo repo.TokenRepository) TokenService {
	return &tokenService{repo: repo}
}

func (t tokenService) CheckToken(tokenVal string) model.Token {
	token := t.repo.FindByToken(tokenVal)
	return token
}