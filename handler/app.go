package handler

import (
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
