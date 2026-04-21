package server

import (
	"github.com/labstack/echo/v4"
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
