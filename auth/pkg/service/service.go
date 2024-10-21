package service

import (
	"github.com/ivnstd/AuthenticationService/auth/models"
	"github.com/ivnstd/AuthenticationService/auth/pkg/repository"
)

type Auth interface {
	CreateUser(user models.User) error
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Auth
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repos.Auth),
	}
}
