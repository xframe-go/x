package liey

import (
	"errors"

	"cnb.cool/liey/liey-go/auth"
)

func RegisterAuth(fn func() auth.Config) {
	cfg := fn()
	rocket.auth = auth.NewManager(cfg)
}

func Auth() *auth.Manager {
	if rocket.auth == nil {
		Logger().Error(errors.New("auth manager is not initialized"))
		return nil
	}
	return rocket.auth
}
