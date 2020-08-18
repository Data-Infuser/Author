package handler

import "gitlab.com/promptech1/infuser-author/app/ctx"

type UserHandler struct {
	Ctx *ctx.Context
}

func NewUserHandler(ctx *ctx.Context) *UserHandler {
	return &UserHandler{
		Ctx: ctx,
	}
}
