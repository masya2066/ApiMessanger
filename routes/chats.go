package routes

import (
	"ApiMessenger/controllers/controllers"
	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine) {
	r.POST("/chat/create", controllers.NewChat)
	r.GET("/chat/list", controllers.ListChat)
	r.GET("/chat/info/:id", controllers.ChatInfo)
	r.DELETE("/chat/delete", controllers.DeleteChat)
}
