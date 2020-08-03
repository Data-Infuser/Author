package service

import (
	"gitlab.com/promptech1/infuser-author/repository"
)

type TokenService interface {

}

type tokenService struct {
	repo repo.TokenRepository
}

func NewTokenService(repo repo.TokenRepository) TokenService {
	return &tokenService{repo: repo}
}
