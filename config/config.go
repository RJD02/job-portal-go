package config

import (
	"RJD02/job-portal/db"
	"log"

	_ "github.com/lib/pq"
)

type Environment string

const (
	Production  Environment = "PRODUCTION"
	Development Environment = "DEVELOPMENT"
)

type Config struct {
	Db               *db.PrismaClient
	JWT_SECRET_KEY   string
	FROM_GMAIL       string
	TO_GMAIL         string
	GMAIL_PASSWORD   string
	ENVIRONMENT      Environment
	ADMIN_SECRET_KEY string
}

var AppConfig Config = Config{}

func (appConfig *Config) Connect(db *db.PrismaClient) {

	if err := db.Prisma.Connect(); err != nil {
		log.Fatal("Error connecting to database", err)
		panic(err)
	}
	appConfig.Db = db
}

func (appConfig *Config) AddSecretKey(key string) {
	if key == "" {
		panic("No SECRET_KEY set")
	}
	appConfig.JWT_SECRET_KEY = key
}

func (appConfig *Config) AddGmailCreds(FROM_GMAIL string, GMAIL_PASSWORD string, TO_GMAIL string) {
	if FROM_GMAIL == "" || GMAIL_PASSWORD == "" || TO_GMAIL == "" {
		panic("GMAIL Credentials not set")
	}
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

func (appConfig *Config) SetAdminKey(key string) {
	if key == "" {
		panic("No ADMIN_KEY set")
	}
	appConfig.ADMIN_SECRET_KEY = key
}
