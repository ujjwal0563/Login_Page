package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"login-api/db"
	"login-api/models"
	"login-api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type SignupRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Phone       string `json:"phone"`
}

func Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if !utils.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not hash password",
		})
		return
	}

	// Handle phone and username fallbacks
	phoneNumber := req.PhoneNumber
	if phoneNumber == "" {
		phoneNumber = req.Phone
	}
	username := req.Username
	if username == "" {
		username = req.Surname
	}

	user := models.User{
		Name:          req.Name,
		Surname:       req.Surname,
		Username:      username,
		Email:         req.Email,
		Password:      hashedPassword,
		PhoneNumber:   phoneNumber,
		IsBlackListed: false,
		CreatedAt:     time.Now(),
	}

	_, err = db.UserCollection.InsertOne(ctx, user)
	if err != nil {
		fmt.Printf("InsertOne error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("User creation failed: %v", err),
		})
		return
	}

	// User successfully created
	log.Printf("User registered successfully: %s\n", user.Email)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
	})
}
