package repo

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/promptech1/infuser-author/model"
)

type AppRepository interface {
	FindByNameSpace(nameSpace string) *model.App
}

type appRepositoryDB struct {
	DB *gorm.DB
}

func NewAppRepository(db *gorm.DB) AppRepository {
	return &appRepositoryDB{DB: db}
}

func (r *appRepositoryDB) FindByNameSpace(nameSpace string) *model.App {
	app := model.App{}

	if r.DB.Debug().Where("name_space = ?", nameSpace).First(&app).RecordNotFound() {
		return nil
	}

	return &app
}