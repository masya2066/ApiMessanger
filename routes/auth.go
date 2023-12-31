package routes

import (
	"ApiMessenger/controllers/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/auth/login", controllers.Login)
	r.POST("/auth/verify", controllers.Verification)
	r.GET("/auth/home", controllers.Home)
	r.GET("/auth/premium", controllers.Prem)
	r.GET("/auth/logout", controllers.Logout)
}
