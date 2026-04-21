package liey

import (
	"cnb.cool/liey/liey-go/contracts"
	"cnb.cool/liey/liey-go/logger"
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
