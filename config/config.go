package config

import (
	"RJD02/job-portal/db"

	_ "github.com/lib/pq"
)

type Environment string

const (
	Production  Environment = "PRODUCTION"
	Development Environment = "DEVELOPMENT"
)

type Config struct {
	Db             *db.PrismaClient
	JWT_SECRET_KEY string
	FROM_GMAIL     string
	TO_GMAIL       string
	GMAIL_PASSWORD string
	ENVIRONMENT    Environment
}

var AppConfig Config = Config{}

func (appConfig *Config) Connect(db *db.PrismaClient) {
	appConfig.Db = db
}

func (appConfig *Config) AddSecretKey(key string) {
	appConfig.JWT_SECRET_KEY = key
}

func (appConfig *Config) AddGmailCreds(FROM_GMAIL string, GMAIL_PASSWORD string, TO_GMAIL string) {
	appConfig.TO_GMAIL = TO_GMAIL
	appConfig.GMAIL_PASSWORD = GMAIL_PASSWORD
	appConfig.FROM_GMAIL = FROM_GMAIL
}

func (appConfig *Config) SetEnv(env string) {
	if env == string(Production) {
		appConfig.ENVIRONMENT = Production
	} else {
		appConfig.ENVIRONMENT = Development
	}
}
