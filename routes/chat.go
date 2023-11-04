package routes

import (
	"ApiMessenger/controllers"
	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine) {
	r.POST("/chat/create", controllers.NewChat)
}
