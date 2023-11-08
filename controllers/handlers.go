package controllers

import (
	"ApiMessenger/models"
	"ApiMessenger/utils"
)

func GetMessagesOfChat(chat models.Chat, count int) ([]models.Message, error) {
	var messages []models.Message
	var er error
	models.DB.Model(&models.Message{}).
		Where("chat_id = ?", chat.ChatId).
		Order("created DESC").
		Limit(count).
		Find(&messages)

	for i := range messages {
		mes, err := utils.DecryptMessage(messages[i].Message, chat.Phrase)
		if err != nil {
			er = err
		}
		messages[i].Message = mes
	}

	if er != nil {
		return messages, er
	} else {
		return messages, nil
	}
}
