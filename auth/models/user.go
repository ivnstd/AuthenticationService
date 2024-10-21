package models

func (User) TableName() string {
	return "users"
}

type User struct {
	ID       int    `json:"id"       gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"     gorm:"not null" binding:"required"`
	Username string `json:"username" gorm:"not null" binding:"required"`
	Password string `json:"password" gorm:"not null" binding:"required"`
}

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
