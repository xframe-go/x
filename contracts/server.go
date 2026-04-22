package contracts

import (
	"github.com/labstack/echo/v4"
)

type Server interface {
	Router
}

type Router interface {
	Get(path string, handler HandlerFunc)
}

type HandlerFunc func(ctx echo.Context) error

type ResourceHandler interface {
	List(ctx echo.Context) error
	BatchList(ctx echo.Context) error
	Create(ctx echo.Context) error
	Show(ctx echo.Context) error
	Update(c echo.Context) error
	Destroy(c echo.Context) error
}
