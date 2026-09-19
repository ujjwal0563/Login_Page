package main

import (
	"log"
	"os"

	"login-api/db"
	"login-api/handlers"
	middlewares "login-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (running in production)")
	}
	db.ConnectMongo()

	r := gin.Default()

    r.Use(cors.Default())

	r.POST("/signup", handlers.Signup)
	r.POST("/register", handlers.Signup)
	r.POST("/login",
		middlewares.LoginLimiter(),
		handlers.Login,
	)
	r.POST("/logout", handlers.Logout)
	r.POST("/forgot-password", handlers.ForgotPassword)
	r.POST("/reset-password", handlers.ResetPassword)

	// Serve frontend UI
	r.StaticFile("/", "./auth-frontend/index.html")
	r.Static("/static", "./auth-frontend")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("\n\n==================================================\n🚀 Server running! Open in Firefox or your browser:\n👉 http://localhost:%s\n==================================================\n\n", port)

	r.Run(":" + port)
}

