package middlewares

import (
	"ApiMessenger/controllers/controllers"
	"ApiMessenger/language"
	"ApiMessenger/utils"

	"github.com/gin-gonic/gin"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")

		if err != nil {
			c.JSON(401, controllers.ErrorMsg(11, language.Language("invalid_login")))
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(cookie)

		if err != nil {
			c.JSON(401, controllers.ErrorMsg(11, language.Language("invalid_login")))
			c.Abort()
			return
		}

		c.Set("role", claims.Role)
		c.Next()
	}
}
