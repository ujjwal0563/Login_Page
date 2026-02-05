package main

import (
	"log"
	"os"

	"login-api/db"
	"login-api/handlers"
	middlewares "login-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (running in production)")
	}
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}

