package db

import (
	"github.com/xframe-go/x/event"
)

type (
	WithConnection interface {
		Connection() string
	}

	Config struct {
		Databases map[string]DriverConf
		EventBus  *event.Bus
	}

	DriverConf struct {
		Driver   string
		Host     string
		Port     uint
		Username string
		Password string
		DB       string
		Charset  string
		Debug    bool
	}
)
