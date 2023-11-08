package models

type RMQMessage struct {
	SessionLost bool   `json:"session_lost"`
	ChatId      string `json:"chat_id"`
	ChatDeleted bool   `json:"chat_deleted"`
	ChatCreated bool   `json:"chat_created"`
}
