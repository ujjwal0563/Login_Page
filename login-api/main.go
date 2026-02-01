package main

import (
	"login-api/db"
	"login-api/handlers"
	middlewares "login-api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectMongo()
	r := gin.Default()
	r.POST("/signup", handlers.Signup)
	r.POST("/login",
		middlewares.LoginLimiter(),
		handlers.Login,
	)
	r.POST("/logout", handlers.Logout)
	r.POST("/forgot-password", handlers.ForgotPassword)
	r.POST("/reset-password", handlers.ResetPassword)

	r.Run(":8080")
}