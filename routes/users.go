package routes

import (
	"ApiMessenger/controllers/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/user/reset/password", controllers.ResetPassword)
	r.GET("/user/info", controllers.UserInfo)
}
