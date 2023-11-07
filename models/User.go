package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Number   string `gorm:"unique" json:"number"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreatedUser struct {
	ID      int    `json:"ID"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Number  string `json:"number"`
	Role    string `json:"role"`
	Created string `json:"CreatedAt"`
	Updated string `json:"UpdatedAt"`
}

type ResetPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
