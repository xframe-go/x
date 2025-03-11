package xhttp

import (
	"github.com/jinzhu/inflection"
	"github.com/labstack/echo/v4"
	"github.com/xframe-go/x/contracts"
)

type echoRouter interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Any(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) []*echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group
	Use(middleware ...echo.MiddlewareFunc)
}

type Router struct {
	engine echoRouter
	cfg    contracts.ServerConfig
	app    contracts.Application
}

func (x *Router) Use(m ...contracts.Middleware) {
	x.engine.Use(wrapMiddlewares(x.cfg, x.app, m...)...)
}

func (x *Router) Get(uri string, handler contracts.Handler) contracts.Router {
	x.engine.GET(uri, wrapHandler(handler, x.cfg, x.app))
	return x
}

func (x *Router) Post(uri string, handler contracts.Handler) contracts.Router {
	x.engine.POST(uri, wrapHandler(handler, x.cfg, x.app))
	return x
}

func (x *Router) Put(uri string, handler contracts.Handler) contracts.Router {
	x.engine.PUT(uri, wrapHandler(handler, x.cfg, x.app))
	return x
}

func (x *Router) Patch(uri string, handler contracts.Handler) contracts.Router {
	x.engine.PATCH(uri, wrapHandler(handler, x.cfg, x.app))
	return x
}

func (x *Router) Delete(uri string, handler contracts.Handler) contracts.Router {
	x.engine.DELETE(uri, wrapHandler(handler, x.cfg, x.app))
	return x
}

func (x *Router) Options(uri string, handler contracts.Handler) contracts.Router {
	x.engine.OPTIONS(uri, wrapHandler(handler, x.cfg, x.app))
	return x
}

func (x *Router) Head(uri string, handler contracts.Handler) contracts.Router {
	x.engine.HEAD(uri, wrapHandler(handler, x.cfg, x.app))
	return x
}

func (x *Router) Connect(uri string, handler contracts.Handler) contracts.Router {
	x.engine.CONNECT(uri, wrapHandler(handler, x.cfg, x.app))
	return x
}

func (x *Router) Trace(uri string, handler contracts.Handler) contracts.Router {
	x.engine.TRACE(uri, wrapHandler(handler, x.cfg, x.app))
	return x
}

func (x *Router) Any(uri string, handler contracts.Handler) contracts.Router {
	x.engine.Any(uri, wrapHandler(handler, x.cfg, x.app))
	return x
}

func (x *Router) Prefix(prefix string, m ...contracts.Middleware) contracts.Router {
	g := x.engine.Group(prefix, wrapMiddlewares(x.cfg, x.app, m...)...)
	return &Router{
		engine: g,
		cfg:    x.cfg,
		app:    x.app,
	}
}

func (x *Router) Resource(resource string, rest contracts.RestHandler) contracts.Router {
	name := inflection.Plural(resource)
	slug := inflection.Singular(resource)
	return x.Get(name, rest.Index).
		Post(name, rest.Store).
		Get(name+":"+slug, rest.Show).
		Get(name+":"+slug, rest.Edit).
		Get(name+":"+slug, rest.Destroy)
}
