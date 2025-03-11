package xhttp

import (
	"fmt"
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/net/xhttp"
	"github.com/xframe-go/x/utils/xmap"
)

type ServerManager struct {
	app     contracts.Application
	cfg     map[string]contracts.ServerConfig
	servers *xmap.Map[string, contracts.XServer]
}

func NewServerManager(app contracts.Application, cfg map[string]contracts.ServerConfig) *ServerManager {
	return &ServerManager{
		app:     app,
		cfg:     cfg,
		servers: xmap.NewMap[string, contracts.XServer](true),
	}
}

func (sm *ServerManager) Server(name string) contracts.XServer {
	srv, _ := sm.servers.GetOrSet(name, func() contracts.XServer {
		srvCfg, ok := sm.cfg[name]
		if !ok {
			sm.app.Log().Fatal(fmt.Sprintf("server '%s' not found", name))
		}

		return xhttp.New(srvCfg, sm.app)
	})
	return srv
}
