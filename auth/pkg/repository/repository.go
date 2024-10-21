package repository

import (
	"github.com/ivnstd/AuthenticationService/auth/models"
	"gorm.io/gorm"
)

type Auth interface {
	CreateUser(user models.User) error
	GetUser(username, password string) (models.User, error)
}

type Repository struct {
	Auth
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Auth: NewAuthDB(db),
	}
}
