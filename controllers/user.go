package controllers

import (
	"ApiMessenger/language"
	"ApiMessenger/models"
	"github.com/gin-gonic/gin"
)

func ResetPassword(c *gin.Context) {

	_, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	var resetPass models.ResetPassword

	if err := c.ShouldBindJSON(&resetPass); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	if resetPass.OldPassword == resetPass.NewPassword {
		c.JSON(400, gin.H{"error": "passwords are same"})
	}

	if len(resetPass.NewPassword) < 8 {
		c.JSON(400, ErrorMsg(16, "Password can't be less 8 characters"))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Password reset was success",
	})

}
