package xdb

import "github.com/xframe-go/x/contracts"

type symbol string

var xdbConfigSymbol = symbol("xdb")

type (
	Config struct {
		Default string

		Connections map[string]contracts.DbDriver
	}
)

func (Config) Symbol() any {
	return xdbConfigSymbol
}

type DatabaseConfig struct {
	Default string
}
