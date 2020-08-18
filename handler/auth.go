package handler

import (
	"gitlab.com/promptech1/infuser-author/app/ctx"
)

type AuthHandler struct {
	Ctx *ctx.Context
}

func NewAuthHandler(ctx *ctx.Context) *AuthHandler {
	return &AuthHandler{Ctx: ctx}
}
