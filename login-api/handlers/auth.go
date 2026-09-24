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

func Login(c *gin.Context) {
	var req LoginRequest
	var user models.User

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.UserCollection.FindOne(
		ctx,
		bson.M{"email": req.Email},
	).Decode(&user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Check if user is blacklisted in database
	if user.IsBlackListed || user.BlackList {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are blacklisted and cannot login",
		})
		return
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Wrong password",
			"action": "forgot-password",
		})
		return
	}

	token, _ := utils.GenerateToken(user.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
