package x

import (
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/logger"
)

func (r *Rocket) getLogger() contracts.Logger {
	if r.logger != nil {
		return r.logger
	}
	r.logger = logger.NewSlog()

	return r.logger
}

func Logger() contracts.Logger {
	return rocket.getLogger()
}
