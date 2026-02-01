package handlers

import (
	"net/http"
	"time"

	"login-api/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ForgotPassword(c *gin.Context) {

	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "If email exists, reset link sent",
		})
		return
	}

	token := uuid.NewString()
	expiry := time.Now().Add(15 * time.Minute)

	err = db.SaveResetToken(user.ID, token, expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server error",
		})
		return
	}

	println("RESET TOKEN:", token)

	c.JSON(http.StatusOK, gin.H{
		"message": "Reset link sent",
	})
}


