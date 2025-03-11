package xhttp

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/xframe-go/x/contracts"
)

type XHttp struct {
	Router
	engine *echo.Echo
	cfg    contracts.ServerConfig
}

func New(cfg contracts.ServerConfig, app contracts.Application) *XHttp {
	if cfg.JsonMarshal == nil {
		cfg.JsonMarshal = json.Marshal
	}

	if cfg.JsonUnmarshal == nil {
		cfg.JsonUnmarshal = json.Unmarshal
	}

	engine := echo.New()

	return &XHttp{
		engine: engine,
		cfg:    cfg,
		Router: Router{
			engine: engine,
			cfg:    cfg,
			app:    app,
		},
	}
}

func (x *XHttp) Start() error {
	// register routes
	if x.cfg.RoutingProviders != nil {
		for _, provider := range x.cfg.RoutingProviders {
			provider(x)
		}
	}

	// register global middleware
	if x.cfg.Middlewares != nil {
		x.Use(x.cfg.Middlewares...)
	}

	addr := fmt.Sprintf("%s:%d", x.cfg.Host, x.cfg.Port)
	//x.engine.HideBanner = true
	//x.engine.HidePort = true
	x.engine.Debug = true
	return x.engine.Start(addr)
}
