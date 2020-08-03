package repo

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/promptech1/infuser-author/model"
	"time"
)

type AppTokenRepository interface {
	Find(appID uint, tokenID uint) *model.AppToken
	FindTodayUsage(appToken *model.AppToken) int
}

type appTokenRepositoryDB struct {
	DB *gorm.DB
}

func (a appTokenRepositoryDB) Find(appID uint, tokenID uint) *model.AppToken {
	appToken := model.AppToken{
		AppID: appID, TokenID: tokenID,
	}

	if a.DB.First(&appToken).RecordNotFound() {
		return nil
	}

	return &appToken
}

func (a appTokenRepositoryDB) FindTodayUsage(appToken *model.AppToken) int {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tomorrow := today.AddDate(0, 0, 1)

	var count int

	a.DB.Where(
		"created_at BETWEEN ? AND ?", today, tomorrow,
	).Find(&appToken).Count(&count)

	return count
}

func NewAppTokenRepository(db *gorm.DB) AppTokenRepository {
	return &appTokenRepositoryDB{DB: db}
}