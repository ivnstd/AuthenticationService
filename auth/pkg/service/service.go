package service

import (
	"github.com/ivnstd/AuthenticationService/auth/models"
	"github.com/ivnstd/AuthenticationService/auth/pkg/repository"
)

type Auth interface {
	CreateUser(user models.User) error
	GetUserByID(id int) (models.User, error)
	GetUserByUsername(username, password string) (models.User, error)

	GenerateAccessToken(userID int) (string, error)
	GenerateRefreshToken(userID int, clientIP string) (string, error)
	ParseAccessToken(accessToken string) (int, error)
	RefreshTokens(accessToken, refreshToken, clientIP string) (string, string, error)
	RevokeRefreshToken(refreshToken string) error
}

type Service struct {
	Auth
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repos.Auth),
	}
}
