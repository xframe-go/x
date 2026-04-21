package handlers

import (
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
}

func NewContext(c echo.Context) *Context {
	return &Context{c}
}

func (ctx *Context) Validated(pointer any) error {
	if err := (&echo.DefaultBinder{}).BindBody(ctx, &pointer); err != nil {
		return err
	}

	return ctx.Validate(pointer)
}
