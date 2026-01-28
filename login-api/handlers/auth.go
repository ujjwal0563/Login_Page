package handlers

import (
	"context"
	"net/http"
	"time"

	"login-api/db"
	"login-api/models"
	"login-api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	var user models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.UserCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		return
	}

	token, _ := utils.GenerateToken(user.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// check if user already exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, _ := db.UserCollection.CountDocuments(ctx, bson.M{
		"email": req.Email,
	})

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already exists",
		})
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not hash password",
		})
		return
	}

	// create user
	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	_, err = db.UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User creation failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
	})
}

