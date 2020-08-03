package service

import (
	"gitlab.com/promptech1/infuser-author/enum"
	repo "gitlab.com/promptech1/infuser-author/repository"
)

type AppTokenService interface {
	CheckAppToken(token string, nameSpace string) enum.AuthCode
}

type appTokenService struct {
	repo repo.AppTokenRepository
	appRepo repo.AppRepository
	tokenRepo repo.TokenRepository
}

func NewAppTokenService(repo repo.AppTokenRepository, appRepo repo.AppRepository, tokenRepo repo.TokenRepository) AppTokenService {
	return &appTokenService{repo: repo, appRepo: appRepo, tokenRepo: tokenRepo}
}

func (s appTokenService) CheckAppToken(token string, nameSpace string) enum.AuthCode{
	app := s.appRepo.FindByNameSpace(nameSpace)
	if app == nil {
		return enum.UNREGISTERED_SERVICE
	}

	t := s.tokenRepo.FindByToken(token)
	if t == nil {
		return enum.UNAUTHORIZED
	}

	appToken := s.repo.Find(app.ID, t.ID)
	if appToken != nil {
		count := s.repo.FindTodayUsage(appToken)
		if int64(app.MaxTraffic) >= int64(count) {
			return enum.VALID
		}

		return enum.LIMIT_EXCEEDED
	}

	return enum.UNAUTHORIZED
}
