package contracts

import (
	"context"
	"gorm.io/gorm"
)

type DbDriver interface {
	Open() (*gorm.DB, error)
}

type DBProvider interface {
	Connection(name ...string) *gorm.DB
	WithContext(ctx context.Context) *gorm.DB
}
