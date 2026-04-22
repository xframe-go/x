package server

import (
	"github.com/labstack/echo/v4"
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/validate"
)

type EchoServer struct {
	*echo.Echo
}

func NewEcho() *EchoServer {
	e := echo.New()

	e.Validator = &validate.FormRequestValidator{}

	return &EchoServer{
		Echo: e,
	}
}

func (*EchoServer) Resource(group *echo.Group, name string, handler contracts.ResourceHandler) {
	group.GET(name, handler.List)
	group.GET(name+"/_batch", handler.BatchList)
	group.GET(name+"/:id", handler.Show)
	group.POST(name, handler.Create)
	group.PUT(name+"/:id", handler.Update)
	group.DELETE(name+"/:id", handler.Destroy)
}
