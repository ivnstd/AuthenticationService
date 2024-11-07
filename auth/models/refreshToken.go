package models

import "time"

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

type RefreshToken struct {
	ID        int       `json:"id"         gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id"    gorm:"not null;index"`
	Token     string    `json:"token"      gorm:"not null;unique"`
	ClientIP  string    `json:"client_ip"  gorm:"not null;type:varchar(45)"`
	ExpiredAt time.Time `json:"expired_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type RefreshInput struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
