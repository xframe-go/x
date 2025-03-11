package xlog

import "github.com/xframe-go/x/contracts"

type symbol string

var loggingSymbol = symbol("logging")

type Config struct {
	// Default Log Channel
	// This option defines the default log channel that is utilized to write
	// messages to your logs. The value provided here should match one of
	// the channels present in the list of "channels" configured below.
	Default string

	// Log Channels
	// Here you may configure the log channels for your application. XFrame
	// utilizes the log/slog logging library, which includes a variety
	// of powerful log handlers and formatters that you're free to use.
	Channels map[string]contracts.Logger
}

func (*Config) Symbol() any {
	return loggingSymbol
}
