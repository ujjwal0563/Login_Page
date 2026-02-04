package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"login-api/db"
	"login-api/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func generateOTP() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "000000"
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 1 minute resend limit
	if !user.OTPLastSentAt.IsZero() && time.Since(user.OTPLastSentAt) < time.Minute {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Wait 1 minute before requesting OTP again",
		})
		return
	}

	// Max 3 OTP per hour
	if user.OTPRequestCount >= 3 && time.Since(user.OTPLastSentAt) < time.Hour {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "OTP limit reached. Try after 1 hour",
		})
		return
	}

	otp := generateOTP()
	hashedOTP, _ := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	expiresAt := time.Now().Add(5 * time.Minute)

	if err := db.SaveOTP(req.Email, string(hashedOTP), expiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save OTP",
		})
		return
	}

	// Send OTP using Brevo REST API
	if err := utils.SendOTPEmailBrevo(req.Email, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send OTP email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to your email (valid for 5 minutes)",
	})
}
