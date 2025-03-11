package contracts

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Response interface {
	Render(ctx echo.Context) error
}

type ResponseWriter interface {
	Json(data any) Response
	String(data string) Response
	Code(code int) ResponseWriter
	NoContent() Response
	Header() http.Header
}
