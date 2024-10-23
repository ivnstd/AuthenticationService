package repository

import (
	"github.com/ivnstd/AuthenticationService/auth/models"
	"gorm.io/gorm"
)

type Auth interface {
	CreateUser(user models.User) error
	GetUserByID(id int) (models.User, error)
	GetUserByUsername(username, password string) (models.User, error)

	SaveRefreshToken(refreshToken models.RefreshToken) error
	GetRefreshToken(refreshToken string) (models.RefreshToken, error)
	DeleteRefreshToken(refreshToken string) error
}

type Repository struct {
	Auth
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Auth: NewAuthDB(db),
	}
}
