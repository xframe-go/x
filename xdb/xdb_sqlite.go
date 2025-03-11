package xdb

import (
	"github.com/xframe-go/x/contracts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type SqliteConfig struct {
	Path string
}

type xDBSqlite struct {
	cfg SqliteConfig
}

func NewSqlite(cfg SqliteConfig) contracts.DbDriver {
	return &xDBSqlite{
		cfg: cfg,
	}
}

func (x *xDBSqlite) Open() (*gorm.DB, error) {
	dir := filepath.Dir(x.cfg.Path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	return gorm.Open(sqlite.Open(x.cfg.Path), &gorm.Config{})
}
