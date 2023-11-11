package controllers

import (
	"ApiMessenger/language"
	"ApiMessenger/models"
	"ApiMessenger/utils"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func NewMessage(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	if cookie == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	parse, err := utils.ParseToken(cookie)

	var message models.Message
	var user models.User
	var chat models.Chat

	models.DB.Model(&models.User{}).Where("number = ?", parse.Subject).First(&user)

	if user.ID == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	err = c.ShouldBindJSON(&message)
	if message.Message == "" || message.ChatId == "" {
		c.JSON(403, ErrorMsg(50, language.Language("invalid_message")))
		return
	}

	models.DB.Model(&models.Chat{}).Where("chat_id = ?", message.ChatId).First(&chat)

	if chat.ID == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	chatUsers := models.UsersOfChat(message.ChatId)
	var existUser bool

	for i := range chatUsers {
		if chatUsers[i] == int(user.ID) {
			existUser = true
		}
	}

	if existUser == true {
		mes, err := utils.EncryptMessage(message.Message, chat.Phrase)
		if err != nil {
			c.JSON(500, ErrorMsg(51, language.Language("send_message_error")))
			return
		}

		message.Message = mes
		message.UserId, message.Owner = int(user.ID), int(user.ID)
		message.Created, message.Update = time.Now().UTC().Format(os.Getenv("DATE_FORMAT")), time.Now().UTC().Format(os.Getenv("DATE_FORMAT"))
		message.MessageId = utils.GenerateId()

		models.DB.Create(message)
		models.DB.Model(&models.Chat{}).Where("chat_id = ?", message.ChatId).Update("updated", message.Created)

		c.JSON(200, gin.H{
			"success": true,
			"message": language.Language("message_sent"),
			"chat":    message.ChatId,
		})
	} else {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}
}

func DeleteMessages(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	if cookie == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	parse, err := utils.ParseToken(cookie)

	var user models.User

	models.DB.Model(&models.User{}).Where("number = ?", parse.Subject).First(&user)

	if user.ID == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	var body models.DeletingMessage
	//var message models.Message
	var chat models.Chat

	_ = c.ShouldBindJSON(&body)

	models.DB.Model(&chat).Where("chat_id = ?", body.ChatId).First(&chat)

	if chat.ID == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	if users := models.UsersOfChat(body.ChatId); len(users) == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

}
