package models

type User struct {
	Id           uint   `json:"id" gorm:"primarykey"`
	Name         string `json:"name" binding:"required,min=4"`
	Email        string `json:"email" gorm:"unique"`
	Password     string `json:"password" binding:"required,min=6"`
	Role         string `json:"role"`
	RefreshToken string `json:"refresh_token"`
}

