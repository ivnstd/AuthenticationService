package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ivnstd/AuthenticationService/auth/config"
	"github.com/ivnstd/AuthenticationService/auth/models"
	"github.com/ivnstd/AuthenticationService/auth/pkg/repository"
	"github.com/sirupsen/logrus"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 15 * 24 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Auth
}

func NewAuthService(repo repository.Auth) *AuthService {
	return &AuthService{repo: repo}
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(config.Config.Salt)))
}

func (s *AuthService) CreateUser(user models.User) error {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUserByUsername(username, password string) (models.User, error) {
	password = generatePasswordHash(password)
	return s.repo.GetUserByUsername(username, password)
}

func (s *AuthService) GetUserByID(id int) (models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *AuthService) GenerateAccessToken(userID int) (string, error) {
	user, err := s.repo.GetUserByID(userID)

	if err != nil {
		logrus.Errorf("Failed to get user: %v", err)
		return "", err
	}

	if user.ID == 0 {
		return "", errors.New("Invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	signedToken, err := token.SignedString([]byte(config.Config.SecretKey))
	if err != nil {
		logrus.Errorf("Failed to sign access token: %v", err)
		return "", err
	}

	return signedToken, nil
}

func (s *AuthService) GenerateRefreshToken(userID int, clientIP string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
	})

	signedToken, err := token.SignedString([]byte(config.Config.SecretKey))
	if err != nil {
		logrus.Errorf("Failed to sign refresh token: %v", err)
		return "", err
	}

	err = s.repo.SaveRefreshToken(models.RefreshToken{
		UserID:    userID,
		Token:     signedToken,
		ClientIP:  clientIP,
		ExpiredAt: time.Now().Add(refreshTokenTTL),
	})
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *AuthService) ParseAccessToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid string method")
		}
		return []byte(config.Config.SecretKey), nil
	})
	if err != nil {
		logrus.Errorf("Failed to parse access token: %v", err)
		return 0, err
	}

	claimsn, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, errors.New("Invalid token or claims")
	}

	return claimsn.UserId, nil
}

func (s *AuthService) RefreshTokens(accessToken, refreshToken, clientIP string) (string, string, error) {
	_, err := s.ParseAccessToken(accessToken)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", errors.New("Invalid access token signature")
		}

		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors == jwt.ValidationErrorExpired {
			logrus.Debug("Access token expired, proceeding with refresh")
		} else {
			logrus.Errorf("Failed to parse access token: %v", err)
			return "", "", errors.New("Invalid access token")
		}
	}

	rt, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil {
		logrus.Errorf("Failed to get refresh token: %v", err)
		return "", "", err
	}

	if rt.ExpiredAt.Before(time.Now()) {
		return "", "", errors.New("Refresh token expired")
	}

	newAccessToken, err := s.GenerateAccessToken(rt.UserID)
	if err != nil {
		logrus.Errorf("Failed to generate new access token: %v", err)
		return "", "", err
	}

	newRefreshToken, err := s.GenerateRefreshToken(rt.UserID, clientIP)
	if err != nil {
		logrus.Errorf("Failed to generate new refresh token: %v", err)
		return "", "", err
	}

	err = s.repo.DeleteRefreshToken(refreshToken)
	if err != nil {
		logrus.Errorf("Failed to delete old refresh token: %v", err)
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) RevokeRefreshToken(refreshToken string) error {
	return s.repo.DeleteRefreshToken(refreshToken)
}
