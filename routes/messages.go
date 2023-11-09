package routes

import (
	"ApiMessenger/controllers/controllers"
	"github.com/gin-gonic/gin"
)

func MessageRoutes(r *gin.Engine) {
	r.POST("/messages/add", controllers.NewMessage)
	r.POST("/messages/delete", controllers.DeleteMessages)
}
