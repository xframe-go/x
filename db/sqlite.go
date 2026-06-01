package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type SqliteDriver struct{}

func (SqliteDriver) Open(conf DriverConf) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s.db", conf.DB)

	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}
