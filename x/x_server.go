package x

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/xframe-go/x/server"
)

func (r *Rocket) getServer() *server.EchoServer {
	if r.server != nil {
		return r.server
	}

	r.server = server.NewEcho()

	r.server.Use(middleware.RequestLogger())
	r.server.Use(middleware.Recover())
	r.server.Use(middleware.CORS())
	return r.server
}

func Server() *server.EchoServer {
	return rocket.getServer()
}
