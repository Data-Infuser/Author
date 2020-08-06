package repo

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/promptech1/infuser-author/model"
	"time"
)

type AppTokenRepository interface {
	Find(appID uint, tokenID uint) *model.AppToken
	FindTodayUsage(appToken *model.AppToken) uint
	Create(appID uint, tokenID uint) *model.AppToken
}

type appTokenRepositoryDB struct {
	DB *gorm.DB
}

func NewAppTokenRepository(db *gorm.DB) AppTokenRepository {
	return &appTokenRepositoryDB{DB: db}
}

func (r appTokenRepositoryDB) Create(appID uint, tokenID uint) *model.AppToken {
	appToken := &model.AppToken{AppID: appID, TokenID: tokenID}
	r.DB.Create(appToken)

	return appToken
}

func (r appTokenRepositoryDB) Find(appID uint, tokenID uint) *model.AppToken {
	appToken := model.AppToken{}

	if r.DB.Where("app_id = ? AND token_id = ?", appID, tokenID).First(&appToken).RecordNotFound() {
		return nil
	}

	return &appToken
}

func (r appTokenRepositoryDB) FindTodayUsage(appToken *model.AppToken) uint {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tomorrow := today.AddDate(0, 0, 1)

	var count uint

	r.DB.Where(
		"created_at BETWEEN ? AND ?", today, tomorrow,
	).Find(&appToken).Count(&count)

	return count
}