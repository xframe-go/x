package app

import "github.com/xframe-go/x/contracts"

type symbol string

var appConfigSymbol = symbol("app")

type Config struct {
	// Application Name
	// This value is the name of your application, which will be used when the
	// framework needs to place the application's name in a notification or
	// other UI elements where an application name needs to be displayed.
	Name string

	// Application Environment
	// This value determines the "environment" your application is currently
	// running in. This may determine how you prefer to configure various
	// services the application utilizes. Set this in your ".env" file.
	Env string

	Servers map[string]contracts.ServerConfig

	Providers []contracts.Provider
}

func (*Config) Symbol() any {
	return appConfigSymbol
}
