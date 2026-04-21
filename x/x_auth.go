package x

import (
	"errors"

	"github.com/xframe-go/x/auth"
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
