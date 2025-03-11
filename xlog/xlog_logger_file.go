package xlog

import (
	"github.com/xframe-go/x/contracts"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const logDir = "storage/logs"

type FileLogger struct {
	*baseLogger
}

type FileLoggerConfig struct {
	// Level reports the minimum record level that will be logged.
	// The handler discards records with lower levels.
	// If Level is nil, the handler assumes LevelInfo.
	// The handler calls Level.Level for each record processed;
	// to adjust the minimum level dynamically, use a LevelVar.
	Level slog.Leveler

	// ReplaceAttr is called to rewrite each non-group attribute before it is logged.
	// The attribute's value has been resolved (see [Value.Resolve]).
	// If ReplaceAttr returns a zero Attr, the attribute is discarded.
	//
	// The built-in attributes with keys "time", "level", "source", and "msg"
	// are passed to this function, except that time is omitted
	// if zero, and source is omitted if AddSource is false.
	//
	// The first argument is a list of currently open groups that contain the
	// Attr. It must not be retained or modified. ReplaceAttr is never called
	// for Group attributes, only their contents. For example, the attribute
	// list
	//
	//     Int("a", 1), Group("g", Int("b", 2)), Int("c", 3)
	//
	// results in consecutive calls to ReplaceAttr with the following arguments:
	//
	//     nil, Int("a", 1)
	//     []string{"g"}, Int("b", 2)
	//     nil, Int("c", 3)
	//
	// ReplaceAttr can be used to change the default keys of the built-in
	// attributes, convert types (for example, to replace a `time.Time` with the
	// integer seconds since the Unix epoch), sanitize personal information, or
	// remove attributes from the output.
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

	Format string

	Writer io.Writer
}

func NewFileLogger(cfg FileLoggerConfig) contracts.Logger {
	var logger *slog.Logger
	if cfg.Format == JSONFormatter {
		logger = slog.New(slog.NewJSONHandler(cfg.Writer, &slog.HandlerOptions{
			AddSource:   false,
			Level:       cfg.Level,
			ReplaceAttr: cfg.ReplaceAttr,
		}))
	}

	if cfg.Format == TextFormatter {
		logger = slog.New(slog.NewTextHandler(cfg.Writer, &slog.HandlerOptions{
			AddSource:   false,
			Level:       cfg.Level,
			ReplaceAttr: cfg.ReplaceAttr,
		}))
	}

	return &FileLogger{
		baseLogger: &baseLogger{
			logger: logger,
		},
	}
}

type FileRotate interface {
	FileName(name string) string
}

type FuncRotate func(name string) string

func (f FuncRotate) FileName(name string) string {
	return f(name)
}

var (
	SingleRotate = FuncRotate(func(name string) string {
		return name
	})

	DailyRotate = FuncRotate(func(name string) string {
		return name + "-" + time.Now().Format(time.DateOnly)
	})
)

type FileWriter struct {
	mu     sync.Mutex
	fs     *os.File
	rotate FileRotate
	name   string

	fileName string
}

type WriterOption func(*FileWriter)

func WithRotate(rotate FileRotate) WriterOption {
	return func(f *FileWriter) {
		f.rotate = rotate
	}
}

func NewFileWriter(name string, options ...WriterOption) io.Writer {
	w := &FileWriter{
		mu:   sync.Mutex{},
		name: name,
	}

	for _, option := range options {
		option(w)
	}

	return w
}

func (fl *FileWriter) Write(p []byte) (n int, err error) {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	fs, err := fl.getFile()
	if err != nil {
		return 0, err
	}
	return fs.Write(p)
}

func (fl *FileWriter) getFile() (*os.File, error) {
	fileName := fl.rotate.FileName(fl.name)
	if fileName == fl.fileName {
		return fl.fs, nil
	}

	if fl.fs != nil {
		_ = fl.fs.Close()
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logFilePath := filepath.Join(logDir, fileName+".log")
	fs, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	fl.fs = fs
	fl.fileName = fileName

	return fl.fs, nil
}
