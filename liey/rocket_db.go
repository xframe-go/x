package liey

import (
	"cnb.cool/liey/liey-go/db"
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
