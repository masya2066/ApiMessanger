package models

import "time"

type UserMessage struct {
	UserId    int       `json:"user_id"`
	ChatId    string    `json:"chat_id"`
	MessageId string    `json:"message_id"`
	Message   string    `json:"message"`
	Changed   bool      `json:"changed"`
	Deleted   bool      `json:"deleted"`
	Created   time.Time `json:"created"`
	Update    time.Time `json:"update"`
}
