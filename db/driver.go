package db

import "gorm.io/gorm"

type Driver interface {
	Open(conf DriverConf) (*gorm.DB, error)
}
