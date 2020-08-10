package repo

import (
	"gitlab.com/promptech1/infuser-author/model"
	"xorm.io/xorm"
)

type AppRepository interface {
	FindByNameSpace(nameSpace string) (*model.App, error)
}

type appRepositoryDB struct {
	DB *xorm.Engine
}

func NewAppRepository(db *xorm.Engine) AppRepository {
	return &appRepositoryDB{DB: db}
}

func (r *appRepositoryDB) FindByNameSpace(nameSpace string) (*model.App, error) {
	app := model.App{NameSpace: nameSpace}

	if _, err := r.DB.Get(&app); err != nil {
		return nil, err
	}

	return &app, nil
}