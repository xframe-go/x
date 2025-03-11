package middlewares

import (
	"github.com/xframe-go/x"
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/frame/xhttp/events"
)

func RequestHandled(req contracts.Request, next func() error) error {
	defer func() {
		p := &events.RequestData{
			Method: req.Method(),
			Path:   req.Path(),
		}
		if err := events.RequestHandledEvent.Emit(p); err != nil {
			x.App().Log().Error(err.Error())
		}
	}()

	if err := next(); err != nil {
		return err
	}

	return nil
}
