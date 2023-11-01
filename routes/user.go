package routes

import (
	"ApiMessenger/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/reset/password", controllers.ResetPassword)
}
