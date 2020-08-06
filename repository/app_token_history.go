package repo

import (
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"gitlab.com/promptech1/infuser-author/model"
)

type AppTokenHistoryRepository interface {
	Create(histories []model.AppTokenHistory)
}

type appTokenHistoryRepositoryDB struct {
	DB *gorm.DB
}

func NewAppTokenHistoryRepository(db *gorm.DB) AppTokenHistoryRepository {
	return &appTokenHistoryRepositoryDB{DB: db}
}

func (r appTokenHistoryRepositoryDB) Create(histories []model.AppTokenHistory) {
	var insertRecords []interface{}

	for _, h := range histories {
		insertRecords = append(insertRecords, h)
	}

	if err := gormbulk.BulkInsert(r.DB, insertRecords, 3000); err != nil {
		glog.Info(err)
	}
}