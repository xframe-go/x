package contracts

type App struct {
	/*
	   |--------------------------------------------------------------------------
	   | Application Name
	   |--------------------------------------------------------------------------
	   |
	   | This value is the name of your application, which will be used when the
	   | framework needs to place the application's name in a notification or
	   | other UI elements where an application name needs to be displayed.
	   |
	*/
	// Application Name
	// This value is the name of your application, which will be used when the
	// framework needs to place the application's name in a notification or
	// other UI elements where an application name needs to be displayed.
	Name string

	Servers map[string]ServerConfig
}

type Application interface {
	Config() Config
	Log() Log
	Server(name ...string) XServer
	DB() DBProvider
}
