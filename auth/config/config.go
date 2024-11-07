package config

import (
	"os"

	"github.com/joho/godotenv"
)

var Config struct {
	Port string

	SecretKey string
	Salt      string

	DB_Host     string
	DB_Port     string
	DB_Username string
	DB_Name     string
	DB_SSLMode  string
	DB_Password string
}

func LoadConfig() error {
	err := godotenv.Load()

	Config.Port = os.Getenv("PORT")

	Config.SecretKey = os.Getenv("SECRET_KEY")
	Config.Salt = os.Getenv("SALT")

	Config.DB_Host = os.Getenv("DB_HOST")
	Config.DB_Port = os.Getenv("DB_PORT")
	Config.DB_Username = os.Getenv("DB_USER")
	Config.DB_Name = os.Getenv("DB_NAME")
	Config.DB_SSLMode = os.Getenv("SSL_MODE")
	Config.DB_Password = os.Getenv("DB_PASSWORD")

	return err
}
