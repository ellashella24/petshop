package config

import (
	"os"
	"sync"

	"petshop/constants"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type AppConfig struct {
	Port     string
	Driver   string
	Name     string
	Address  string
	DB_Port  string
	Username string
	Password string
}

var lock = &sync.Mutex{}
var appConfig *AppConfig

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}

	return appConfig
}

func initConfig() *AppConfig {
	err := godotenv.Load(".env")
	if err != nil {
		log.Infof("can't read file env: %s", err)
	}

	var defaultConfig AppConfig
	defaultConfig.Port = os.Getenv("PORT")
	defaultConfig.Driver = os.Getenv("DRIVER")
	defaultConfig.Name = os.Getenv("DB_NAME")
	defaultConfig.Address = os.Getenv("ADDRESS")
	defaultConfig.DB_Port = os.Getenv("DB_PORT")
	defaultConfig.Username = os.Getenv("DB_USERNAME")
	defaultConfig.Password = os.Getenv("DB_PASSWORD")

	constants.SecretKey = os.Getenv("SecretKey")
	constants.CallbackToken = os.Getenv("CallbackToken")
	constants.XendToken = os.Getenv("XendToken")
	constants.FTP_ADDR = os.Getenv("FTP_ADDR")
	constants.FTP_USERNAME = os.Getenv("FTP_USERNAME")
	constants.FTP_PASSWORD = os.Getenv("FTP_PASSWORD")
	constants.CONFIG_SMTP_HOST = os.Getenv("CONFIG_SMTP_HOST")
	constants.CONFIG_SMTP_PORT = os.Getenv("CONFIG_SMTP_PORT")
	constants.CONFIG_SENDER_NAME = os.Getenv("CONFIG_SENDER_NAME")
	constants.CONFIG_AUTH_EMAIL = os.Getenv("CONFIG_AUTH_EMAIL")
	constants.CONFIG_AUTH_PASSWORD = os.Getenv("CONFIG_AUTH_PASSWORD")

	return &defaultConfig
}
