package handler

import (
	"time"

	"github.com/thoas/go-funk"
	"gitlab.com/promptech1/infuser-author/app/ctx"
	"gitlab.com/promptech1/infuser-author/model"
)

type AppHandler struct {
	Ctx *ctx.Context
}

func NewAppHandler(ctx *ctx.Context) *AppHandler {
	return &AppHandler{Ctx: ctx}
}

func (h *AppHandler) Create(app *model.App) error {
	session := h.Ctx.Orm.NewSession()
	session.Begin()

	if _, err := h.Ctx.Orm.Insert(app); err != nil {
		session.Rollback()
		return err
	}

	if _, err := h.Ctx.Orm.Insert(app.Traffics); err != nil {
		session.Rollback()
		return err
	}

	if _, err := h.Ctx.Orm.Insert(app.Operations); err != nil {
		session.Rollback()
		return err
	}

	session.Commit()

	return nil
}

func (h *AppHandler) Update(app *model.App) error {
	session := h.Ctx.Orm.NewSession()
	session.Begin()

	if _, err := h.Ctx.Orm.ID(app.Id).Update(app); err != nil {
		return err
	}

	originOperations, err := model.FindOperationsByApp(h.Ctx.Orm, app.Id)
	if err != nil {
		return err
	}
	var originIds []uint
	for _, operation := range originOperations {
		originIds = append(originIds, operation.Id)
		operation.DelRedis(h.Ctx.RedisDB) // 수정이 발생할 수 있는 operation에 대한 redis 삭제 처리
	}

	var newIds []uint
	for _, operation := range app.Operations {
		newIds = append(newIds, operation.Id)
	}

	//기존 데이터와 공통되는(Update 대상) ID 추출
	updatedIds := funk.Join(originIds, newIds, funk.InnerJoin)
	for _, id := range updatedIds.([]uint) {
		idx := funk.IndexOf(newIds, id)
		operation := app.Operations[idx]

		err := operation.Update(h.Ctx.Orm)
		if err != nil {
			h.Ctx.Logger.Info(err)
			session.Rollback()
			return err
		}
	}

	//기존 데이터와 차이가 있는(Delete / Insert) ID 추출
	deleteIds, insertIds := funk.Difference(originIds, newIds)
	for _, id := range deleteIds.([]uint) {
		idx := funk.IndexOf(originIds, id)
		delOperation := originOperations[idx]
		err := delOperation.Delete(h.Ctx.Orm)
		if err != nil {
			h.Ctx.Logger.Info(err)
			session.Rollback()
			return err
		}
	}

	operations := []model.Operation{}
	if len(insertIds.([]uint)) > 0 {
		for _, id := range insertIds.([]uint) {
			idx := funk.IndexOf(newIds, id)
			operations = append(operations, app.Operations[idx])
		}
		if _, err := h.Ctx.Orm.Insert(operations); err != nil {
			session.Rollback()
			return err
		}
	}

	traffics, err := model.FindTrafficsByApp(h.Ctx.Orm, app.Id)
	h.Ctx.Logger.Debug(traffics)
	for _, traffic := range traffics {
		traffic.DelRedis(h.Ctx.RedisDB)
		err := traffic.Delete(h.Ctx.Orm)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	if _, err := h.Ctx.Orm.Insert(app.Traffics); err != nil {
		session.Rollback()
		return err
	}

	session.Commit()

	return nil
}

func (h *AppHandler) Destroy(appId uint) error {
	session := h.Ctx.Orm.NewSession()
	session.Begin()

	// 1. Operation 삭제 처리
	var operations []model.Operation
	if err := h.Ctx.Orm.Where("app_id = ?", appId).Find(&operations); err != nil {
		session.Rollback()
		return err
	}
	// 1-1. DB 삭제 처리
	operationSql := "UPDATE operation SET deleted_at = ? WHERE app_id = ? AND deleted_at IS NULL"
	if _, err := h.Ctx.Orm.Exec(operationSql, time.Now(), appId); err != nil {
		session.Rollback()
		return err
	}
	// 1-2. Redis 삭제 처리
	for _, operation := range operations {
		operation.DelRedis(h.Ctx.RedisDB)
	}

	// 2. Traffic 삭제 처리
	var traffics []model.Traffic
	if err := h.Ctx.Orm.Where("app_id = ?", appId).Find(&traffics); err != nil {
		session.Rollback()
		return err
	}
	// 2-1. DB 삭제 처리
	trafficSql := "DELETE FROM traffic WHERE app_id = ?"
	if _, err := h.Ctx.Orm.Exec(trafficSql, appId); err != nil {
		session.Rollback()
		return err
	}

	// 2-2. Redis 삭제 처리
	for _, traffic := range traffics {
		traffic.DelRedis(h.Ctx.RedisDB)
	}

	app := &model.App{Id: appId}

	if err := app.FindApp(h.Ctx.Orm); err != nil {
		session.Rollback()
		return err
	}

	app.DelRedis(h.Ctx.RedisDB)
	app.Delete(h.Ctx.Orm)

	session.Commit()

	return nil
}
