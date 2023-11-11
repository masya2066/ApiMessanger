package controllers

import (
	"ApiMessenger/consumers"
	"ApiMessenger/language"
	"ApiMessenger/models"
	"ApiMessenger/utils"
	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	var user models.User

	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	claims, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	models.DB.Model(&models.User{}).Where("number = ?", claims.Subject).First(&user)

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
