package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/spf13/cast"
)

func JWTMiddleware(authManager *Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			token, err := getToken(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": err.Error(),
				})
			}

			claims, err := authManager.ParseToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "invalid or expired token",
				})
			}

			ctx := c.Request().Context()

			ctx = SetEnterpriseID(ctx, claims.EnterpriseID)

			ctx = SetUserID(ctx, cast.ToUint64(claims.ID))

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func getToken(c *echo.Context) (string, error) {
	cookie, err := c.Cookie("token")
	if err == nil && cookie != nil {
		token := cookie.Value
		if len(token) > 0 {
			return token, nil
		}
	}

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	return parts[1], nil
}
