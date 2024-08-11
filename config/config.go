package config

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Config struct {
	Db *sql.DB
}

var AppConfig Config = Config{}

func (appConfig *Config) Connect(db *sql.DB) {
	appConfig.Db = db
}
