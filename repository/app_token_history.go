package repo

import (
	"gitlab.com/promptech1/infuser-author/model"
	"xorm.io/xorm"
)

type AppTokenHistoryRepository interface {
	Create(histories []model.AppTokenHistory)
}

type appTokenHistoryRepositoryDB struct {
	DB *xorm.Engine
}

func NewAppTokenHistoryRepository(db *xorm.Engine) AppTokenHistoryRepository {
	return &appTokenHistoryRepositoryDB{DB: db}
}

func (r appTokenHistoryRepositoryDB) Create(histories []model.AppTokenHistory) {
	var insertRecords []interface{}

	for _, h := range histories {
		insertRecords = append(insertRecords, h)
	}

	// TODO
	//if err := gormbulk.BulkInsert(r.DB, insertRecords, 3000); err != nil {
	//	glog.Info(err)
	//}
}