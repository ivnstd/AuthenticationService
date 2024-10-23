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
		logrus.Errorf("Failed to create user: %v", err)
		return err
	}

	logrus.Debug("User created successfully")
	return nil
}

func (r *AuthDB) GetUserByID(id int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error

	if err != nil {
		logrus.Errorf("Failed to get user: %v", err)
		return user, err
	}

	logrus.Debug("User fetched successfully")
	return user, nil
}

func (r *AuthDB) GetUserByUsername(username, password string) (models.User, error) {
	var user models.User
	err := r.db.Where("username = ? AND password = ?", username, password).First(&user).Error

	if err != nil {
		logrus.Errorf("Failed to get user: %v", err)
		return user, err
	}

	logrus.Debug("User fetched successfully")
	return user, nil
}

func (r *AuthDB) SaveRefreshToken(refreshToken models.RefreshToken) error {
	err := r.db.Create(&refreshToken).Error

	if err != nil {
		logrus.Errorf("Failed to save refresh token: %v", err)
		return err
	}

	logrus.Debug("Refresh token saved successfully")
	return nil
}

func (r *AuthDB) GetRefreshToken(refreshToken string) (models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.Preload("User").Where("token = ?", refreshToken).First(&token).Error

	if err != nil {
		logrus.Errorf("Failed to get refresh token: %v", err)
		return token, err
	}

	logrus.Debug("Refresh token fetched successfully")
	return token, nil
}

func (r *AuthDB) DeleteRefreshToken(refreshToken string) error {
	err := r.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{}).Error

	if err != nil {
		logrus.Errorf("Failed to delete refresh token: %v", err)
		return err
	}

	logrus.Debug("Refresh token deleted successfully")
	return nil
}
