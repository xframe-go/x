package db

import (
	"fmt"

	"gorm.io/gorm"
)

type DB struct {
	cfg Config

	connections map[string]*gorm.DB

	drivers map[string]Driver
}

func New(cfg Config) *DB {
	return &DB{
		cfg: cfg,
		drivers: map[string]Driver{
			"mysql":    &MysqlDriver{},
			"postgres": &PostgresDriver{},
		},
		connections: make(map[string]*gorm.DB),
	}
}

func (db *DB) Connect() error {
	for name := range db.cfg.Databases {
		conn, err := db.doConnect(name)
		if err != nil {
			return err
		}
		db.connections[name] = conn
	}
	return nil
}

func (db *DB) doConnect(name string) (*gorm.DB, error) {
	cfg := db.cfg.Databases

	conf, ok := cfg[name]
	if !ok {
		return nil, fmt.Errorf("unknown db name: %s", name)
	}

	driver, ok := db.drivers[conf.Driver]
	if !ok {
		return nil, fmt.Errorf("unknown driver: %s", conf.Driver)
	}

	instance, err := driver.Open(conf)
	if err != nil {
		return nil, err
	}

	if conf.Debug {
		instance = instance.Debug()
	}

	return instance, nil
}

func (db *DB) DB(conn ...string) (*gorm.DB, error) {
	if db.connections == nil || len(db.connections) == 0 {
		return nil, fmt.Errorf("database connection is empty")
	}

	name := "default"
	if len(conn) > 0 {
		name = conn[0]
	}

	instance, ok := db.connections[name]
	if !ok {
		return nil, fmt.Errorf("not found database connection: %s", name)
	}
	return instance, nil
}
