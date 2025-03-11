package app

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/xframe-go/x/cmd/x"
	"github.com/xframe-go/x/configs"
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/frame/xdb"
	"github.com/xframe-go/x/frame/xhttp"
	"github.com/xframe-go/x/utils/singleton"
	"github.com/xframe-go/x/xevent"
	"github.com/xframe-go/x/xlog"
	"log"
)

type Application struct {
	xConfig  *singleton.Singleton[*configs.Config]
	xLog     *singleton.Singleton[*xlog.Logging]
	xServers *singleton.Singleton[*xhttp.ServerManager]
	xCommand *singleton.Singleton[*x.Command]
	xEvent   *singleton.Singleton[contracts.EventBus]
	xDB      *singleton.Singleton[contracts.DBProvider]
}

func New() *Application {
	var envFile string
	flag.StringVar(&envFile, "env", ".env", "env file")
	flag.Parse()

	// 加载环境变量
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal(err)
	}

	cfg := singleton.New(configs.New)
	cfg.Get()

	logging := singleton.New(func() *xlog.Logging {
		return xlog.New(cfg.Get())
	})

	app := &Application{
		xConfig: cfg,
		xLog:    logging,
		xEvent:  singleton.New(xevent.New),
	}

	app.xServers = singleton.New(func() *xhttp.ServerManager {
		config := app.AppConfig()
		return xhttp.NewServerManager(app, config.Servers)
	})

	app.xDB = singleton.New(func() contracts.DBProvider {
		return xdb.New(app)
	})

	app.xCommand = singleton.New(func() *x.Command {
		return x.New(app)
	})

	return app
}

func (app *Application) AppConfig() *Config {
	cfg, ok := app.xConfig.Get().Get(appConfigSymbol)
	if !ok {
		app.Log().Fatal("missing application configuration")
	}
	return cfg.(*Config)
}

func (app *Application) Config() contracts.Config {
	return app.xConfig.Get()
}

func (app *Application) Log() contracts.Log {
	return app.xLog.Get()
}

func (app *Application) Server(names ...string) contracts.XServer {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}
	return app.xServers.Get().Server(name)
}

func (app *Application) Event() contracts.EventBus {
	return app.xEvent.Get()
}

func (app *Application) DB() contracts.DBProvider {
	return app.xDB.Get()
}

func (app *Application) Execute() error {
	app.boot()

	return app.xCommand.Get().Execute()
}

func (app *Application) boot() {
	// init logger
	app.Log()

	conf, ok := app.Config().Get(appConfigSymbol)
	if !ok {
		app.Log().Fatal("missing application configuration")
	}

	for i := range conf.(*Config).Providers {
		provider := conf.(*Config).Providers[i]
		provider.Register(app)
	}
}
