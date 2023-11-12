package controllers

import (
	"ApiMessenger/consumers"
	"ApiMessenger/language"
	"ApiMessenger/middlewares"
	"ApiMessenger/models"
	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	var user models.User

	isAuth, parse := middlewares.IsAuthorized(c)

	if !isAuth || parse.Subject == "" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	models.DB.Model(&models.User{}).Where("number = ?", parse.Subject).First(&user)

	if user.ID == 0 {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	c.JSON(200, gin.H{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"number":  "+" + user.Number,
		"created": user.Created,
		"updated": user.Updated,
	})

	consumers.SendJSON(models.RMQMessage{
		SessionLost: false,
	})
}
