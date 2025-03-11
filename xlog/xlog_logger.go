package xlog

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/xframe-go/x/contracts"
)

type baseLogger struct {
	ctx    context.Context
	logger *slog.Logger
}

func (b baseLogger) Context(ctx context.Context) contracts.Logger {
	return &baseLogger{
		logger: b.logger,
		ctx:    ctx,
	}
}

func (b baseLogger) Info(msg string, args ...slog.Attr) {
	b.log(slog.LevelInfo, msg, args...)
}

func (b baseLogger) Debug(msg string, args ...slog.Attr) {
	b.log(slog.LevelDebug, msg, args...)
}

func (b baseLogger) Warn(msg string, args ...slog.Attr) {
	b.log(slog.LevelWarn, msg, args...)
}

func (b baseLogger) Error(msg string, args ...slog.Attr) {
	b.log(slog.LevelError, msg, args...)
}

func (b baseLogger) Fatal(msg string, args ...slog.Attr) {
	b.log(slog.LevelError, msg, args...)
	os.Exit(1)
}

func (b baseLogger) log(level slog.Level, msg string, args ...slog.Attr) {
	var ctx context.Context
	skip := 3
	if b.ctx != nil {
		ctx = b.ctx
	} else {
		ctx = context.Background()
		skip = 4
	}

	if !b.logger.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	// skip 调用栈:
	// 1. runtime.Callers
	// 2. this function (log)
	// 3. wrapper functions (Info/Debug/Warn/Error/Fatal)
	// 4. actual caller
	runtime.Callers(skip, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.AddAttrs(args...)
	file, line := runtime.FuncForPC(pc).FileLine(pc)
	r.AddAttrs(slog.String("caller", file+":"+strconv.Itoa(line-1)))
	_ = b.logger.Handler().Handle(ctx, r)
}
