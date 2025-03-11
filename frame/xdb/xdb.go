package xdb

import (
	"context"
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/utils/xmap"
	"gorm.io/gorm"
	"log/slog"
)

type XDB struct {
	app  contracts.Application
	conn *xmap.Map[string, *gorm.DB]
}

func New(app contracts.Application) contracts.DBProvider {
	return &XDB{
		app:  app,
		conn: xmap.NewMap[string, *gorm.DB](true),
	}
}

func (xdb *XDB) Connection(name ...string) *gorm.DB {
	return xdb.getConnection(name...)
}

func (xdb *XDB) getConnection(name ...string) *gorm.DB {
	cfg, ok := xdb.app.Config().Get(xdbConfigSymbol)
	if !ok {
		xdb.app.Log().Fatal("missing db config")
	}

	conf, ok := cfg.(*Config)
	if !ok {
		xdb.app.Log().Fatal("invalid db config")
	}

	driver := conf.Default
	if len(name) > 0 && len(name[0]) > 0 {
		driver = name[0]
	}

	if conf.Connections == nil || len(conf.Connections) == 0 {
		xdb.app.Log().Fatal("missing connection config")
	}

	conn, ok := conf.Connections[driver]
	if !ok {
		xdb.app.Log().Fatal("missing connection '" + driver + "' config")
	}

	tx, ok := xdb.conn.Get(driver)
	if ok {
		return tx
	}

	tx, err := conn.Open()
	if err != nil {
		xdb.app.Log().Fatal("failed to open connection,", slog.String("connection", driver), slog.String("err", err.Error()))
	}
	xdb.conn.Set(driver, tx)
	return tx
}

func (xdb *XDB) WithContext(ctx context.Context) *gorm.DB {
	return xdb.getConnection().WithContext(ctx)
}
