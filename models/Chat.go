package models

type Chat struct {
	ID      int    `gorm:"unique" json:"id"`
	ChatId  string `gorm:"unique" json:"chat_Id"`
	Name    string `json:"name"`
	Public  bool   `json:"public"`
	Owner   uint   `json:"owner"`
	Phrase  string `json:"phrase"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type ChatMembers struct {
	UserId      int    `json:"user_id"`
	ChatId      string `json:"chat_id"`
	Owner       bool   `json:"owner"`
	Role        string `json:"role"`
	DateCreated string `json:"created"`
	DateUpdated string `json:"updated"`
}

type ChatInfo struct {
	Name    string `json:"name"`
	ChatId  string `json:"chat_id"`
	Members []int  `json:"members"`
	Owner   int    `json:"owner"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

func UsersOfChat(chatsId string) []int {
	var chat []ChatMembers
	var users []int
	DB.Model(ChatMembers{}).Where("chat_id = ?", chatsId).Find(&chat)
	for _, chatMember := range chat {
		users = append(users, chatMember.UserId)
	}
	return users
}
