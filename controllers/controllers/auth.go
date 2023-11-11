package controllers

import (
	"ApiMessenger/language"
	"ApiMessenger/models"
	"ApiMessenger/utils"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
)

var jwtKey = []byte("$2a$14$cTtWCEgT1Vl2Q1orRe1GRObekuWfmW3DEhsf.I9kgoG45SkMg/Y.2")

func ErrorMsg(code int, mes string) map[string]any {
	return gin.H{
		"success": false,
		"code":    code,
		"message": mes,
	}
}

func Login(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if num := utils.CheckDigits(user.Number); num != true {
		c.JSON(400, ErrorMsg(4, language.Language("invalid_number")))
		return
	}

	var existingUser models.User

	models.DB.Where("number = ?", user.Number).First(&existingUser)

	if existingUser.ID == 0 || existingUser.Number == "" {
		models.DB.Model(&models.User{}).Create(&models.User{
			Name:    user.Name,
			Email:   user.Email,
			Number:  user.Number,
			Role:    "user",
			Created: time.Now().UTC().Format(os.Getenv("DATE_FORMAT")),
			Updated: time.Now().UTC().Format(os.Getenv("DATE_FORMAT")),
			Deleted: false,
		})
	}

	var userInfo models.CreatedUser

	result := models.DB.Where("number = ?", user.Number).First(&user)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	models.DB.Model(&user).Where("number = ?", user.Number).First(&userInfo)

	access := models.CheckAccessToSendSms(user.Number)

	if access == false {
		c.JSON(200, gin.H{
			"user":     userInfo,
			"sms_send": access,
		})
		return
	}

	send := models.DB.Create(models.SmsCode{
		Number:   user.Number,
		Code:     utils.GenerateVerify(),
		Sent:     true,
		Attempts: 0,
		SentTime: time.Now().UTC().Format(os.Getenv("DATE_FORMAT")),
		Created:  time.Now().UTC().Format(os.Getenv("DATE_FORMAT")),
	})

	if send.Error != nil {
		c.JSON(400, gin.H{"error": send.Error})
		return
	}
	c.JSON(200, gin.H{
		"user":     userInfo,
		"sms_send": true,
	})
}

func Verification(c *gin.Context) {
	var user models.User
	var sms models.SmsCode
	var body models.SmsBody

	_ = c.ShouldBindJSON(&body)

	models.DB.Model(&user).Where("number = ?", body.Number).First(&user)

	if body.Number == "" || utils.CheckDigits(body.Number) == false {
		c.JSON(400, ErrorMsg(52, language.Language("invalid_number")))
		return
	}

	models.DB.Model(&models.SmsCode{}).Where("number = ?", user.Number).First(&sms)

	if sms.Number == "" || sms.Created == "" {
		c.JSON(400, ErrorMsg(57, language.Language("invalid_login")))
		return
	}

	count, verdict := models.AttemptSubmitSms(body.Number, body.Code)

	strToIntAttempts, err := strconv.Atoi(os.Getenv("VERIFICATION_ATTEMPTS"))
	if err != nil {
		panic(err)
	}

	attempts := strToIntAttempts - count

	if !verdict {
		if attempts == 0 {
			c.JSON(400, ErrorMsg(56, language.Language("input_code_after")))
			return
		}
		c.JSON(400, ErrorMsg(55, language.Language("incorrect_pin")+strconv.Itoa(attempts)))
		return
	}

	remove := models.DB.Model(&sms).Where("number = ?", body.Number).Delete(&sms)
	if remove.Error != nil {
		panic(remove.Error)
		return
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("JSON marshaling error:", err)
		panic(err)
		return
	}

	bytes := []byte(jsonData)
	decodeError := json.Unmarshal(bytes, &user)
	if decodeError != nil {
		fmt.Println("JSON decoding error:", decodeError)
		return
	}

	TokenLife := os.Getenv("TOKEN_LIFE_TIME")

	life, err := time.ParseDuration(TokenLife)
	if err != nil {
		panic(err.Error())
	}

	expirationTime := time.Now().UTC().Add(life)

	claims := &models.Claims{
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Number,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, ErrorMsg(-1, language.Language("fail_generate_token")))
		return
	}

	models.DB.Model(&user).Where("number = ?", body.Number).Update("active", true)

	models.DB.Model(&user).Where("number = ?", body.Number).First(&user.Active)

	fmt.Println("activate:", user.Active)

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"success": true,
		"user":    user,
		"token":   tokenString,
	})

}

func Home(c *gin.Context) {

	cookie, err := c.Cookie("token")

	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	claims, err := utils.ParseToken(cookie)

	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	if claims.Role != "user" && claims.Role != "admin" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	c.JSON(200, gin.H{"success": "home page", "role": claims.Role})
}

func Prem(c *gin.Context) {

	cookie, err := c.Cookie("token")

	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	claims, err := utils.ParseToken(cookie)

	if err != nil {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	if claims.Role != "admin" {
		c.JSON(401, ErrorMsg(11, language.Language("invalid_login")))
		return
	}

	c.JSON(200, gin.H{"success": "premium page", "role": claims.Role})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": true})
}
