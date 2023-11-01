package controllers

import (
	"ApiMessenger/language"
	"ApiMessenger/models"
	"ApiMessenger/utils"
	"github.com/gin-gonic/gin"
)

func ResetPassword(c *gin.Context) {

	claims, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	email, err := utils.ParseToken(claims)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	var reset models.ResetPassword

	if err := c.ShouldBindJSON(&reset); err != nil {
		ErrorMsg(403, err.Error())
	}

	var usr models.User

	_ = models.DB.Model(&models.User{}).Where("email = ?", email.Subject).First(&usr)

	newPass, err := utils.GenerateHashPassword(reset.NewPassword)
	if err != nil {
		ErrorMsg(403, err.Error())
		return
	}

	compare := utils.CompareHashPassword(reset.OldPassword, usr.Password)
	if compare == false {
		c.JSON(400, ErrorMsg(400, "Old password is not correct!"))
		return
	}

	if len(reset.NewPassword) < 8 {
		c.JSON(400, ErrorMsg(16, "Password can't be less 8 characters"))
		return
	}

	if reset.OldPassword == reset.NewPassword {
		c.JSON(400, ErrorMsg(15, "Password can't be the same!"))
		return
	}

	models.DB.Model(&models.User{}).Where("email = ?", email.Subject).Update("password", newPass)
	c.JSON(200, gin.H{
		"success": true,
		"message": "Password reset was success",
	})

}
