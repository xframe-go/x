package xhttp

import (
	"github.com/labstack/echo/v4"
	"github.com/xframe-go/x/contracts"
	"net/http"
)

type JsonResponse struct {
	cfg    contracts.ServerConfig
	header http.Header
	data   any
	code   int
}

func (r *JsonResponse) Render(ctx echo.Context) error {
	for key := range r.header {
		ctx.Response().Header().Set(key, r.header.Get(key))
	}
	return ctx.JSON(r.code, r.data)
}

type XResponse struct {
	req    *Request
	cfg    contracts.ServerConfig
	code   int
	header http.Header
}

func (x *XResponse) Code(code int) contracts.ResponseWriter {
	x.code = code
	return x
}

func (x *XResponse) Header() http.Header {
	return x.header
}

func (x *XResponse) Json(data any) contracts.Response {
	return &JsonResponse{
		data:   data,
		cfg:    x.cfg,
		code:   x.code,
		header: x.header,
	}
}

func (x *XResponse) NoContent() contracts.Response {
	return &EmptyResponse{
		header: x.header,
	}
}

func (x *XResponse) String(data string) contracts.Response {
	return &StringResponse{
		data:   data,
		code:   x.code,
		header: x.header,
	}
}

type StringResponse struct {
	code   int
	data   string
	header http.Header
}

func (r *StringResponse) Render(ctx echo.Context) error {
	for key := range r.header {
		ctx.Response().Header().Set(key, r.header.Get(key))
	}
	return ctx.String(r.code, r.data)
}

type EmptyResponse struct {
	header http.Header
}

func (r *EmptyResponse) Render(ctx echo.Context) error {
	for key := range r.header {
		ctx.Response().Header().Set(key, r.header.Get(key))
	}
	return ctx.NoContent(http.StatusNoContent)
}
