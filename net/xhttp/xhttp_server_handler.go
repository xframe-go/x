package xhttp

import (
	"github.com/labstack/echo/v4"
	"github.com/xframe-go/x/contracts"
)

func wrapHandler(handler contracts.Handler, cfg contracts.ServerConfig, app contracts.Application) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			req = createRequestFromEchoContext(ctx, cfg, app)

			resp = handler(req)
		)

		if resp == nil {
			return nil
		}

		return resp.Render(ctx)
	}
}

func wrapMiddleware(middleware contracts.Middleware, cfg contracts.ServerConfig, app contracts.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := createRequestFromEchoContext(ctx, cfg, app)
			return middleware(req, func() error {
				return next(ctx)
			})
		}
	}
}

func wrapMiddlewares(cfg contracts.ServerConfig, app contracts.Application, ms ...contracts.Middleware) (middlewares []echo.MiddlewareFunc) {
	for i := range ms {
		m := ms[i]
		middlewares = append(middlewares, wrapMiddleware(m, cfg, app))
	}
	return
}

func EchoMiddleware(middlewareFunc echo.MiddlewareFunc) contracts.Middleware {
	return func(request contracts.Request, next func() error) error {
		req := request.(*Request)
		fn := middlewareFunc(func(c echo.Context) error {
			return next()
		})
		return fn(req.Context)
	}
}
