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
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	var reset models.ResetPassword

	if err := c.ShouldBindJSON(&reset); err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	var usr models.User

	_ = models.DB.Model(&models.User{}).Where("email = ?", email.Subject).First(&usr)

	newPass, err := utils.GenerateHashPassword(reset.NewPassword)
	if err != nil {
		c.JSON(403, ErrorMsg(-1, err.Error()))
		return
	}

	compare := utils.CompareHashPassword(reset.OldPassword, usr.Password)
	if compare == false {
		c.JSON(400, ErrorMsg(400, language.Language("incorrect_old_pass")))
		return
	}

	if len(reset.NewPassword) < 8 {
		c.JSON(400, ErrorMsg(16, language.Language("short_pass")))
		return
	}

	if reset.OldPassword == reset.NewPassword {
		c.JSON(400, ErrorMsg(15, language.Language("same_passwords")))
		return
	}

	models.DB.Model(&models.User{}).Where("email = ?", email.Subject).Update("password", newPass)
	c.JSON(200, gin.H{
		"success": true,
		"message": language.Language("success_reset_pass"),
	})

}

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

	models.DB.Model(&models.User{}).Where("email = ?", claims.Subject).First(&user)

	c.JSON(200, gin.H{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"number": "+" + user.Number,
		"create": user.CreatedAt,
		"update": user.UpdatedAt,
	})

}
