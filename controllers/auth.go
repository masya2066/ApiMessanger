package controllers

import (
	"ApiMessenger/language"
	"ApiMessenger/models"
	"ApiMessenger/utils"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

var jwtKey = []byte("my_secret_key")

func ErrorMsg(code int, mes string) map[string]any {
	return gin.H{
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

	var existingUser models.User

	models.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(400, ErrorMsg(13, language.Language("user_not_exist")))
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		c.JSON(400, ErrorMsg(14, language.Language("invalid_password")))
		return
	}

	var userInfo models.CreatedUser

	result := models.DB.Where("email = ?", user.Email).First(&user)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("JSON marshaling error:", err)
	}

	bytes := []byte(jsonData)
	decodeError := json.Unmarshal(bytes, &userInfo)
	if decodeError != nil {
		fmt.Println("JSON decoding error:", decodeError)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(500, ErrorMsg(-1, language.Language("fail_generate_token")))
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"user":  userInfo,
		"token": tokenString,
	})
}

func Signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, ErrorMsg(-1, err.Error()))
		return
	}

	user.Role = "user"

	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.JSON(403, ErrorMsg(14, language.Language("invalid_reg_data")))
		fmt.Println(&user)
		return
	}

	if user.Number == "" {
		c.JSON(403, ErrorMsg(15, "Number is empty"))
		return
	}

	var existingUser models.User

	models.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID != 0 {
		c.JSON(400, ErrorMsg(12, language.Language("account_already_exist")))
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		c.JSON(500, ErrorMsg(-1, language.Language("fail_generate_pass_hash")))
		return
	}

	models.DB.Create(&user)

	c.JSON(200, gin.H{
		"success": true,
		"message": language.Language("success_signup"),
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
