package routes

import (
	"login-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Public routes (NO JWT)
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)

	// OTP based password reset (NO JWT)
	r.POST("/forgot-password", handlers.ForgotPassword)
	r.POST("/reset-password", handlers.ResetPassword)
}