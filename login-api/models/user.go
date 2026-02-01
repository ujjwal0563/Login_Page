package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Surname          string             `bson:"surname" json:"surname"`
	Email            string             `bson:"email" json:"email"`
	Password         string             `bson:"password" json:"-"`
	PhoneNumber      string             `bson:"phone_number" json:"phone_number"`
	ResetToken       string             `bson:"reset_token,omitempty"`
	ResetTokenExpiry *time.Time         `bson:"reset_token_expiry,omitempty"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
}