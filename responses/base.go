package responses

import (
	"errors"
	"net/http"

	"cnb.cool/liey/liey-go/validate"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Base struct {
}

func (Base) Success(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusOK, data)
}

func (Base) Created(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusCreated, data)
}

func (Base) Failed(ctx echo.Context, err error) error {
	var (
		msg  = err.Error()
		code = http.StatusBadRequest
	)

	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		msg = validate.Translate(errs[0])
		code = http.StatusUnprocessableEntity
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		code = http.StatusNotFound
	}

	return ctx.JSON(code, map[string]string{
		"message": msg,
	})
}

func (Base) Empty(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}

type SSE struct {
	*echo.Response
}

func NewSSE(c echo.Context) *SSE {
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	return &SSE{
		Response: w,
	}
}

func (s *SSE) Send(data []byte) error {
	event := Event{
		Data: data,
	}
	if err := event.MarshalTo(s.Response); err != nil {
		return err
	}

	s.Response.Flush()
	return nil
}

func (Base) SSE(ctx echo.Context) *SSE {
	return NewSSE(ctx)
}
