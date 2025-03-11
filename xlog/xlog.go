package xlog

import (
	"context"
	"github.com/xframe-go/x/contracts"
	"log"
	"log/slog"
)

type Logging struct {
	*Logger

	cfg *Config
}

func New(app contracts.Config) *Logging {
	cfg, ok := app.Get(loggingSymbol)
	if !ok {
		log.Fatal("[xlog] logger config not found")
	}
	conf, ok := cfg.(*Config)
	if !ok {
		log.Fatal("[xlog] invalid logger config")
	}

	return &Logging{
		cfg:    conf,
		Logger: NewLogger(getLogger(conf, conf.Default)),
	}
}

func (l *Logging) Channel(name string) contracts.Logger {
	return NewLogger(getLogger(l.cfg, name))
}

func getLogger(cfg *Config, name string) contracts.Logger {
	channels := cfg.Channels
	if channels == nil {
		log.Fatal("logging [channels] can not be nil")
	}

	channel, ok := channels[name]
	if !ok {
		log.Fatal("logging [channel] can not be found in channels [", name, "]")
	}

	return channel
}

type Logger struct {
	logger contracts.Logger
}

func NewLogger(logger contracts.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l Logger) Context(ctx context.Context) contracts.Logger {
	return l.logger.Context(ctx)
}

func (l Logger) Info(msg string, args ...slog.Attr) {
	l.logger.Info(msg, args...)
}

func (l Logger) Fatal(msg string, args ...slog.Attr) {
	l.logger.Fatal(msg, args...)
}

func (l Logger) Debug(msg string, args ...slog.Attr) {
	l.logger.Debug(msg, args...)
}

func (l Logger) Warn(msg string, args ...slog.Attr) {
	l.logger.Warn(msg, args...)
}

func (l Logger) Error(msg string, args ...slog.Attr) {
	l.logger.Error(msg, args...)
}
