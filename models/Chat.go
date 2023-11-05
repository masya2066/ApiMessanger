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
	UserId      int       `json:"user_id"`
	ChatId      string    `json:"chat_id"`
	Owner       int       `json:"owner"`
	Role        string    `json:"role"`
	DateCreated time.Time `json:"created"`
	DateUpdated time.Time `json:"updated"`
}

func UsersOfChat(chatId string) []int {
	var chat []ChatMembers
	var users []int
	DB.Model(ChatMembers{}).Where("chat_id = ?", chatId).Find(&chat)
	for _, chatMember := range chat {
		users = append(users, chatMember.UserId)
	}
	return users
}
