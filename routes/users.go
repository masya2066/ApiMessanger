package routes

import (
	"ApiMessenger/controllers/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/user/info", controllers.UserInfo)
}
