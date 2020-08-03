package service

import (
	"gitlab.com/promptech1/infuser-author/enum"
	repo "gitlab.com/promptech1/infuser-author/repository"
)

type AppTokenService interface {
	CheckAppToken(token string, nameSpace string) enum.AuthCode
	Regist(token string, nameSpace string)
}

type appTokenService struct {
	repo repo.AppTokenRepository
	appRepo repo.AppRepository
	tokenRepo repo.TokenRepository
}

func NewAppTokenService(repo repo.AppTokenRepository, appRepo repo.AppRepository, tokenRepo repo.TokenRepository) AppTokenService {
	return &appTokenService{repo: repo, appRepo: appRepo, tokenRepo: tokenRepo}
}

func (s appTokenService) Regist(token string, nameSpace string) {
	app := s.appRepo.FindByNameSpace(nameSpace)
	t := s.tokenRepo.FindByToken(token)
	if t == nil {
		t = s.tokenRepo.Create(token)
	}

	s.repo.Create(app.ID, t.ID)
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
			// TODO: API 인증 횟수 추가 처리 필요함
			return enum.VALID
		}

		return enum.LIMIT_EXCEEDED
	}

	return enum.UNAUTHORIZED
}
