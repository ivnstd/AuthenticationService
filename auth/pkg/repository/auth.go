package repository

import (
	"github.com/ivnstd/AuthenticationService/auth/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthDB struct {
	db *gorm.DB
}

func NewAuthDB(db *gorm.DB) *AuthDB {
	return &AuthDB{db: db}
}

func (r *AuthDB) CreateUser(user models.User) error {
	err := r.db.Create(&user).Error

	if err != nil {
		logrus.Errorf("Failed to execute query: %v", err)
		return err
	}

	logrus.Debug("Query executed successfully, song creates")
	return nil
}

func (r *AuthDB) GetUser(username, password string) (models.User, error) {
	var user models.User
	err := r.db.Where("username = ? AND password = ?", username, password).First(&user).Error

	if err != nil {
		logrus.Errorf("Failed to execute query: %v", err)
		return user, err
	}

	logrus.Debug("Query executed successfully, song creates")
	return user, nil
}
