package xdb

import (
	"fmt"
	"github.com/xframe-go/x/contracts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlConfig struct {
	User      string `json:"user,omitempty"`
	Password  string `json:"password,omitempty"`
	Host      string `json:"host,omitempty"`
	Port      int    `json:"port,omitempty"`
	Database  string `json:"database,omitempty"`
	Charset   string `json:"charset,omitempty"`
	Collation string `json:"collation,omitempty"`
}

func (cfg MysqlConfig) dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset, cfg.Collation,
	)
}

type xDBMySql struct {
	cfg MysqlConfig
}

func NewMySql(cfg MysqlConfig) contracts.DbDriver {
	return &xDBMySql{
		cfg: cfg,
	}
}

func (x *xDBMySql) Open() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(x.cfg.dsn()), &gorm.Config{})
}
