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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
 
	//  Find user by reset token
	var user struct {
		ID               interface{} `bson:"_id"`
		ResetTokenExpiry time.Time   `bson:"reset_token_expiry"`
		Password         string      `bson:"password"`
	}

	err := db.UserCollection.FindOne(
		ctx,
		bson.M{"reset_token": req.Token},
	).Decode(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid token",
		})
		return
	}

   // Check token expiry
	if time.Now().After(user.ResetTokenExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token expired",
		})
		return
	}

	//  Hash new password (reuse utils)
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not hash password",
		})
		return
	}

	//  Update password & clear token
	_, err = db.UserCollection.UpdateByID(
		ctx,
		user.ID,
		bson.M{
			"$set": bson.M{
				"password": hashedPassword,
			},
			"$unset": bson.M{
				"reset_token":        "",
				"reset_token_expiry": "",
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
