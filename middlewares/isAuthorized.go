package middlewares

import (
	"ApiMessenger/models"
	"ApiMessenger/utils"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"time"
)

func IsAuthorized(c *gin.Context) (bool, *models.Claims) {
	header := c.GetHeader("Authorization")

	if header == "" {
		return false, &models.Claims{}
	}

	parse, err := utils.ParseToken(header)
	if err != nil {
		log.Fatal(err)
		return false, &models.Claims{}
	}

	user := models.User{}

	var token models.UserToken

	models.DB.Model(&token).Where("token = ?", header).First(&token)

	lifeTime, err := strconv.Atoi(os.Getenv("TOKEN_LIFE_TIME"))
	if err != nil {
		log.Fatal("Error token life time:", err)
	}

	lifeTime = -lifeTime

	now := time.Now().UTC().Add(time.Second * time.Duration(lifeTime)).Format(os.Getenv("DATE_FORMAT"))

	if now >= token.Created {
		models.DB.Model(&token).Where("token = ?", header).Delete(&token)
		return false, &models.Claims{}
	}

	models.DB.Model(&user).Where("number = ?", parse.Subject).First(&user)

	if user.ID == 0 || user.Number == "" {
		models.DB.Model(&token).Where("token = ?", header).Delete(&token)
		return false, &models.Claims{}
	}

	return true, parse
}
