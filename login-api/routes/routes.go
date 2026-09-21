package routes

import (
	"login-api/db"
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
	r.POST("/admin/blacklist", handlers.BlacklistUser(db.UserCollection))
	r.POST("/admin/unblacklist", handlers.UnblacklistUser(db.UserCollection))
	r.GET("/admin/users", handlers.ListUsers(db.UserCollection))
	r.PUT("/admin/users/:id/blacklist", handlers.BlacklistUser(db.UserCollection))
	r.PUT("/admin/users/:id/unblacklist", handlers.UnblacklistUser(db.UserCollection))
	r.GET("/admin/users/:id/blacklist", handlers.GetBlacklistStatus(db.UserCollection))
}
