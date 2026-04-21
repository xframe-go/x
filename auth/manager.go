package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Manager struct {
	config Config
}

func NewManager(config Config) *Manager {
	return &Manager{
		config: config,
	}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (m *Manager) GenerateToken(userID string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.config.Expiration) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.Secret))
}

func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.config.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func (m *Manager) GetUserId(c echo.Context) (string, error) {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return "", nil
	}
	return userID, nil
}

func (m *Manager) getToken(c echo.Context) (string, error) {
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
