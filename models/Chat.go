package models

import (
	"time"
)

type Chat struct {
	ID      int       `gorm:"unique" json:"id"`
	ChatId  string    `json:"chat_id"`
	Name    string    `json:"name"`
	Public  bool      `json:"public"`
	Owner   uint      `json:"owner"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type ChatMembers struct {
	UserId    int       `json:"user_id"`
	ChatId    string    `json:"chat_id"`
	Owner     bool      `json:"owner"`
	Role      string    `json:"role"`
	DateAdded time.Time `json:"date_added"`
}
