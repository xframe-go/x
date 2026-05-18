package x

import (
	"github.com/xframe-go/x/db"
	"github.com/xframe-go/x/event"
	"gorm.io/gorm"
)

func RegisterDB(fn func() db.Config) {
	cfg := fn()

	if cfg.Databases == nil {
		return
	}

	rocket.db = db.New(cfg)

	if err := rocket.db.Connect(); err != nil {
		panic(err)
	}

	plugin := event.NewPlugin(rocket.bus, event.GormPluginConfig{
		PublishCreated: true,
		PublishUpdated: true,
		PublishDeleted: true,
		Prefix:         "liey",
	})
	for name := range cfg.Databases {
		instance, err := rocket.db.DB(name)
		if err != nil {
			return
		}
		if err = instance.Use(plugin); err != nil {
			panic(err)
		}
	}
}

func DB(conn ...string) *gorm.DB {
	instance, err := rocket.db.DB(conn...)
	if err != nil {
		Logger().Error(err)
		return nil
	}

	return instance
}

func Model[T any](tx ...*gorm.DB) gorm.Interface[T] {
	if len(tx) == 0 {
		var m *T

		name := "default"
		if conn, ok := any(m).(db.WithConnection); ok {
			name = conn.Connection()
		}

		return gorm.G[T](DB(name))
	}

	return gorm.G[T](tx[0])
}
