package database

import (
	"github.com/gentwolf-shen/gohelper/config"
)

var (
	drivers map[string]*Database
)

func init() {
	drivers = make(map[string]*Database, 2)
}

func LoadFromConfig(configs map[string]config.DbConfig) error {
	for name, cfg := range configs {
		if err := AddDriver(name, cfg); err != nil {
			CloseAll()

			return err
		}
	}
	return nil
}

func AddDriver(name string, cfg config.DbConfig) error {
	db := New()
	if err := db.Open(cfg.Type, cfg.Dsn, cfg.MaxOpenConnections, cfg.MaxIdleConnections); err != nil {
		return err
	}
	drivers[name] = db

	return nil
}

func Driver(name string) *Database {
	return drivers[name]
}

func Close(name string) {
	db, ok := drivers[name]
	if ok {
		_ = db.Close()
		delete(drivers, name)
	}
}

func CloseAll() {
	for name, db := range drivers {
		_ = db.Close()
		delete(drivers, name)
	}
}
