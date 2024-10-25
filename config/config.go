package config

import (
	"RJD02/job-portal/db"

	_ "github.com/lib/pq"
)

type Config struct {
	Db             *db.PrismaClient
	JWT_SECRET_KEY string
}

var AppConfig Config = Config{}

func (appConfig *Config) Connect(db *db.PrismaClient) {
	appConfig.Db = db
}

func (appConfig *Config) AddSecretKey(key string) {
	appConfig.JWT_SECRET_KEY = key
}
