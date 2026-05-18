package x

import (
	"github.com/labstack/echo/v5/middleware"
	"github.com/xframe-go/x/server"
)

func Server() *server.EchoServer {
	return rocket.server
}

func RegisterServer(fn func() server.Config) {
	cfg := fn()
	rocket.server = server.NewEcho(cfg)
	rocket.server.Use(middleware.RequestLogger())
	rocket.server.Use(middleware.Recover())
	rocket.server.Use(middleware.CORS("*"))
}
