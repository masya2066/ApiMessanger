package controllers

import (
	"ApiMessenger/language"
	"ApiMessenger/models"
	"ApiMessenger/utils"
	"github.com/gin-gonic/gin"
)

func NewChat(c *gin.Context) {
	var chat models.Chat
	var user models.User

	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	parse, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	models.DB.Model(&models.User{}).Where("email = ?", parse.Subject).First(&user)

	err = c.ShouldBindJSON(&chat)

	if chat.Name == "" {
		c.JSON(403, ErrorMsg(20, language.Language("invalid_chat_name")))
		return
	}

	code := utils.GenerateId()
	chat.Owner = user.ID
	chat.ChatId = code
	models.DB.Create(&chat)

	c.JSON(200, gin.H{
		"success": true,
		"chat_id": code,
		"owner":   user.ID,
		"name":    chat.Name,
	})
}

func DeleteChat(c *gin.Context) {
	var body models.Chat
	var chat models.Chat
	var user models.User

	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	parse, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	err = c.ShouldBindJSON(&body)
	if body.ChatId == "" {
		c.JSON(400, ErrorMsg(21, language.Language("invalid_chat_id")))
		return
	}

	models.DB.Model(&user).Where("email = ?", parse.Subject).First(&user)
	models.DB.Model(&chat).Where("chat_id = ?", body.ChatId).First(&chat)

	if chat.ChatId == "" {
		c.JSON(400, ErrorMsg(25, language.Language("delete_chat_impossible")))
		return
	}

	if chat.Owner != user.ID {
		c.JSON(400, ErrorMsg(25, language.Language("delete_chat_impossible")))
		return
	}

	models.DB.Delete(chat)
	c.JSON(200, gin.H{
		"success": true,
		"message": language.Language("delete_chat_successful"),
	})

}

func ListChat(c *gin.Context) {

}
