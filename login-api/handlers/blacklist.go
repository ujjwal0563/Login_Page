package handlers

import (
	"context"
	"net/http"
	"time"

	"login-api/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// buildUserFilter constructs a BSON filter to search by ObjectID hex, email, username, or phone.
func buildUserFilter(identifier string) bson.M {
	if objID, err := primitive.ObjectIDFromHex(identifier); err == nil {
		return bson.M{
			"$or": []bson.M{
				{"_id": objID},
				{"email": identifier},
				{"username": identifier},
				{"phone_number": identifier},
			},
		}
	}
	return bson.M{
		"$or": []bson.M{
			{"email": identifier},
			{"username": identifier},
			{"phone_number": identifier},
		},
	}
}

// BlacklistUser returns a gin.HandlerFunc that blacklists a user.
func BlacklistUser(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		col := collection
		if col == nil {
			col = db.UserCollection
		}
		if col == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database connection not initialized",
			})
			return
		}

		var req struct {
			Email       string `json:"email"`
			Username    string `json:"username"`
			PhoneNumber string `json:"phone_number"`
		}
		_ = c.ShouldBindJSON(&req)

		target := c.Param("id")
		if target == "" {
			if req.Email != "" {
				target = req.Email
			} else if req.Username != "" {
				target = req.Username
			} else if req.PhoneNumber != "" {
				target = req.PhoneNumber
			}
		}

		if target == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email or user identifier is required",
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		filter := buildUserFilter(target)
		update := bson.M{
			"$set": bson.M{
				"is_blacklisted": true,
				"blacklist":      true,
			},
		}

		res, err := col.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to blacklist user: " + err.Error(),
			})
			return
		}

		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "User blacklisted successfully",
			"is_blacklisted": true,
		})
	}
}

// UnblacklistUser returns a gin.HandlerFunc that removes a user from the blacklist.
func UnblacklistUser(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		col := collection
		if col == nil {
			col = db.UserCollection
		}
		if col == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database connection not initialized",
			})
			return
		}

		var req struct {
			Email       string `json:"email"`
			Username    string `json:"username"`
			PhoneNumber string `json:"phone_number"`
		}
		_ = c.ShouldBindJSON(&req)

		target := c.Param("id")
		if target == "" {
			if req.Email != "" {
				target = req.Email
			} else if req.Username != "" {
				target = req.Username
			} else if req.PhoneNumber != "" {
				target = req.PhoneNumber
			}
		}

		if target == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email or user identifier is required",
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		filter := buildUserFilter(target)
		update := bson.M{
			"$set": bson.M{
				"is_blacklisted": false,
				"blacklist":      false,
			},
		}

		res, err := col.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to unblacklist user: " + err.Error(),
			})
			return
		}

		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "User unblacklisted successfully",
			"is_blacklisted": false,
		})
	}
}

// GetBlacklistStatus returns a gin.HandlerFunc that checks if a user is blacklisted.
func GetBlacklistStatus(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		col := collection
		if col == nil {
			col = db.UserCollection
		}
		if col == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database connection not initialized",
			})
			return
		}

		userID := c.Param("id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User identifier (:id) is required",
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		filter := buildUserFilter(userID)
		var user struct {
			ID            primitive.ObjectID `bson:"_id" json:"id"`
			Email         string             `bson:"email" json:"email"`
			Username      string             `bson:"username" json:"username"`
			IsBlackListed bool               `bson:"is_blacklisted" json:"is_blacklisted"`
			BlackList     bool               `bson:"blacklist" json:"blacklist"`
		}

		err := col.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		isBlacklisted := user.IsBlackListed || user.BlackList

		c.JSON(http.StatusOK, gin.H{
			"id":             user.ID.Hex(),
			"email":          user.Email,
			"username":       user.Username,
			"is_blacklisted": isBlacklisted,
		})
	}
}

// ListUsers returns a gin.HandlerFunc that lists all registered users with their blacklist status.
func ListUsers(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		col := collection
		if col == nil {
			col = db.UserCollection
		}
		if col == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database connection not initialized",
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		cursor, err := col.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch users: " + err.Error(),
			})
			return
		}
		defer cursor.Close(ctx)

		type UserResponse struct {
			ID            primitive.ObjectID `bson:"_id" json:"id"`
			Name          string             `bson:"name" json:"name"`
			Surname       string             `bson:"surname,omitempty" json:"surname,omitempty"`
			Username      string             `bson:"username,omitempty" json:"username,omitempty"`
			Email         string             `bson:"email" json:"email"`
			PhoneNumber   string             `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
			IsBlackListed bool               `bson:"is_blacklisted" json:"is_blacklisted"`
			BlackList     bool               `bson:"blacklist" json:"blacklist"`
			CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
		}

		var users []UserResponse
		if err := cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to decode users",
			})
			return
		}

		for i := range users {
			if users[i].BlackList {
				users[i].IsBlackListed = true
			}
		}

		c.JSON(http.StatusOK, users)
	}
}

