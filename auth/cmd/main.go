package main

import (
	db_server "github.com/ivnstd/AuthenticationService/auth"
	"github.com/ivnstd/AuthenticationService/auth/config"
	"github.com/ivnstd/AuthenticationService/auth/pkg/handler"
	"github.com/ivnstd/AuthenticationService/auth/pkg/repository"
	"github.com/ivnstd/AuthenticationService/auth/pkg/service"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		logrus.Fatalf("Error loading env variables: %s", err.Error())
	}

	logrus.Info("Starting server...")

	db, err := repository.NewDB(repository.Config{
		Host:     config.Config.DB_Host,
		Port:     config.Config.DB_Port,
		Username: config.Config.DB_Username,
		DBName:   config.Config.DB_Name,
		SSLMode:  config.Config.DB_SSLMode,
		Password: config.Config.DB_Password,
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize db: %s", err.Error())
	}
	logrus.Info("Database connection established")

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(db_server.Server)
	if err := srv.Run(config.Config.Port, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error occured while running http server: %s", err.Error())
	}
	logrus.Info("http server successfully launched")
}
