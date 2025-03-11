package xdb

import (
	"context"
	"gorm.io/gorm"
)

type gormDB struct {
	tx *gorm.DB
}

func (db *gormDB) Ping() error {
	s, err := db.tx.DB()
	if err != nil {
		return err
	}
	return s.Ping()
}

func (db *gormDB) WithContext(ctx context.Context) *gorm.DB {
	return db.tx.WithContext(ctx)
}
