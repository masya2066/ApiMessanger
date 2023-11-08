package routes

import (
	"ApiMessenger/controllers/controllers"
	"github.com/gin-gonic/gin"
)

func MessageRoutes(r *gin.Engine) {
	r.POST("/messages/add", controllers.NewMessage)
}
