package server

import (
	"cnb.cool/liey/liey-go/validate"
	"github.com/labstack/echo/v4"
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
