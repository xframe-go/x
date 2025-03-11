package contracts

import (
	"context"
	"log/slog"
)

type Log interface {
	Logger
	Channel(name string) Logger
}

type Logger interface {
	Context(ctx context.Context) Logger
	Info(msg string, args ...slog.Attr)
	Fatal(msg string, args ...slog.Attr)
	Debug(msg string, args ...slog.Attr)
	Warn(msg string, args ...slog.Attr)
	Error(msg string, args ...slog.Attr)
}
