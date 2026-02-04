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
   // Fields for password reset and OTP
	ResetOTPHash     string             `bson:"reset_otp_hash,omitempty"`
	OTPExpiresAt     time.Time          `bson:"otp_expires_at,omitempty"`
	OTPLastSentAt    time.Time          `bson:"otp_last_sent_at,omitempty"`
	OTPRequestCount  int                `bson:"otp_request_count,omitempty"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
}