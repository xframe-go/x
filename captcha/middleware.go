package captcha

import (
	"net/http"

	gocap "github.com/ackcoder/go-cap"
	"github.com/labstack/echo/v4"
)

type Body struct {
	Captcha string `json:"captcha"`
}

func Middleware(cap *gocap.Cap) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			captcha := c.Request().Header.Get("X-Secret")
			if len(captcha) == 0 {
				return c.JSON(http.StatusForbidden, echo.Map{"message": "服务器繁忙，等会儿再试试！"})
			}

			ok := cap.ValidateToken(c.Request().Context(), captcha)
			if !ok {
				return c.JSON(http.StatusForbidden, echo.Map{"message": "服务器繁忙，等会儿再试试！"})
			}
			return next(c)
		}
	}
}
