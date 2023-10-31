package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreatedUser struct {
	ID      int    `json:"ID"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Created string `json:"CreatedAt"`
	Updated string `json:"UpdatedAt"`
}
