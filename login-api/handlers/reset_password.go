package handlers

import (
	"net/http"
	"time"

	"login-api/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(c *gin.Context) {
	var req struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}

	//  Request validation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	//  Password strength check
	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters long",
		})
		return
	}

	//  Get user using DB helper
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	//  OTP requested check
	if user.ResetOTPHash == "" || user.OTPExpiresAt.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP not requested",
		})
		return
	}

	//  OTP expiry check
	if time.Now().After(user.OTPExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP expired",
		})
		return
	}

	//  OTP hash verification
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.ResetOTPHash),
		[]byte(req.OTP),
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid OTP",
		})
		return
	}

	//  Hash new password
	hashedPwd, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	//  Update password
	_, err = db.UserCollection.UpdateOne(
		c,
		map[string]interface{}{"email": req.Email},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"password": hashedPwd,
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Password update failed",
		})
		return
	}

	//  Clear OTP data (one-time use)
	_ = db.ClearOTP(req.Email)

	//  Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successful",
	})
}
