package repo

import (
	"time"

	"gitlab.com/promptech1/infuser-author/model"
	"xorm.io/xorm"
)

type AppTokenRepository interface {
	Find(appID uint, tokenID uint) (*model.AppToken, error)
	FindTodayUsage(appToken *model.AppToken) uint
	Create(appID uint, tokenID uint) *model.AppToken
}

type appTokenRepositoryDB struct {
	DB *xorm.Engine
}

func NewAppTokenRepository(db *xorm.Engine) AppTokenRepository {
	return &appTokenRepositoryDB{DB: db}
}

func (r appTokenRepositoryDB) Create(appID uint, tokenID uint) *model.AppToken {
	appToken := &model.AppToken{AppId: appID, TokenId: tokenID}
	r.DB.Insert(appToken)

	return appToken
}

func (r appTokenRepositoryDB) Find(appID uint, tokenID uint) (*model.AppToken, error) {
	appToken := model.AppToken{
		AppId: appID,
		TokenId: tokenID,
	}

	if _, err := r.DB.Get(&appToken); err != nil {
		return nil, err
	}

	return &appToken, nil
}

func (r appTokenRepositoryDB) FindTodayUsage(appToken *model.AppToken) uint {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tomorrow := today.AddDate(0, 0, 1)
	_ = tomorrow

	var count uint
	_ = count

	//r.DB.Where(
	//	"created_at BETWEEN ? AND ?", today, tomorrow,
	//).Find(&appToken).Count(&count)

	// TODO
	return 0
}