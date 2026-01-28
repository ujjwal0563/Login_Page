package main

import (
	"login-api/db"
	"login-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectMongo()

	r := gin.Default()
	r.POST("/signup", handlers.Signup)

	r.POST("/login", handlers.Login)

	r.Run(":8080")
}