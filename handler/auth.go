package handler

import (
	"gitlab.com/promptech1/infuser-author/app/ctx"
	"gitlab.com/promptech1/infuser-author/model"
)

type AuthHandler struct {
	Ctx *ctx.Context
}

func NewAuthHandler(ctx *ctx.Context) *AuthHandler {
	return &AuthHandler{Ctx: ctx}
}

func (h *AuthHandler) Login(user *model.User) error {
	return nil
}
