package routes

import (
	"login-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/login", handlers.Login)
	r.POST("/signup", handlers.Signup)
	r.POST("/forgot-password", handlers.ForgotPassword)
}