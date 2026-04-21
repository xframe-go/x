package logger

import (
	"log/slog"
	"os"

	"cnb.cool/liey/liey-go/contracts"
)

type Slog struct {
	*slog.Logger
}

func NewSlog() contracts.Logger {
	return &Slog{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (s Slog) Info(msg string) {
	s.Logger.Info(msg)
}

func (s Slog) Error(err error) {
	if err != nil {
		s.Logger.Error(err.Error())
	}
}
