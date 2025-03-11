package xlog

import (
	"github.com/xframe-go/x/contracts"
	"log/slog"
)

type DiscardLogger struct {
	*baseLogger
}

func NewDiscardLogger() contracts.Logger {
	return &DiscardLogger{
		baseLogger: &baseLogger{
			logger: slog.New(slog.DiscardHandler),
		},
	}
}
