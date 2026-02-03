package handlers

import (
	"context"
	"net/http"
	"time"

	"login-api/db"
	"login-api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)
func ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req.Token = strings.TrimSpace(req.Token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user struct {
		ID               interface{} `bson:"_id"`
		ResetToken       string      `bson:"reset_token"`
		ResetTokenExpiry time.Time   `bson:"reset_token_expiry"`
	}

	err := db.UserCollection.FindOne(
		ctx,
		bson.M{"reset_token": req.Token},
	).Decode(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Reset token not found or already used",
		})
		return
	}

	if user.ResetTokenExpiry.IsZero() || time.Now().After(user.ResetTokenExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Reset token expired",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not hash password",
		})
		return
	}

	_, err = db.UserCollection.UpdateByID(
		ctx,
		user.ID,
		bson.M{
			"$set": bson.M{"password": hashedPassword},
			"$unset": bson.M{
				"reset_token":         "",
				"reset_token_expiry":  "",
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Password update failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successful",
	})
}
