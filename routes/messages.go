package routes

import (
	"ApiMessenger/controllers"
	"github.com/gin-gonic/gin"
)

func MessageRoutes(r *gin.Engine) {
	r.POST("/messages/add", controllers.NewMessage)
}
