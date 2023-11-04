package controllers

import (
	"ApiMessenger/language"
	"ApiMessenger/models"
	"ApiMessenger/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

func NewChat(c *gin.Context) {
	var chat models.Chat

	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
	}

	parse, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
	}

	fmt.Println(parse.Subject)

	err = c.ShouldBindJSON(&chat)
	if err != nil {
	}

	fmt.Println(chat.Name)
}
