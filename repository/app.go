package repo

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/promptech1/infuser-author/model"
)

type AppRepository interface {
	FindByNameSpace(nameSpace string) (*model.App, error)
}

type appRepositoryDB struct {
	DB *gorm.DB
}

func NewAppRepository(db *gorm.DB) AppRepository {
	return &appRepositoryDB{DB: db}
}

func (r *appRepositoryDB) FindByNameSpace(nameSpace string) (*model.App, error) {
	app := model.App{}

	if err := r.DB.Where("name_space = ?", nameSpace).First(&app).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return &app, nil
}