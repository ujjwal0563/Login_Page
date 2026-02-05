package db

import (
	"context"
	"time"

	"login-api/models"

	"go.mongodb.org/mongo-driver/bson"
)

//  Get user by email
// Login, forgot-password
func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := UserCollection.FindOne(
		context.TODO(),
		bson.M{"email": email},
	).Decode(&user)

	return user, err
}

// Save OTP hash + expiry + resend tracking
func SaveOTP(
	email string,
	otpHash string,
	expiry time.Time,
) error {
	_, err := UserCollection.UpdateOne(
		context.TODO(),
		bson.M{"email": email},
		bson.M{
			"$set": bson.M{
				"reset_otp_hash":   otpHash,
				"otp_expires_at":   expiry,
				"otp_last_sent_at": time.Now(),
			},
			"$inc": bson.M{
				"otp_request_count": 1,
			},
		},
	)

	return err
}

// Clear OTP data after successful reset
func ClearOTP(email string) error {
	_, err := UserCollection.UpdateOne(
		context.TODO(),
		bson.M{"email": email},
		bson.M{
			"$unset": bson.M{
				"reset_otp_hash":   "",
				"otp_expires_at":   "",
				"otp_last_sent_at": "",
			},
			"$set": bson.M{
				"otp_request_count": 0,
			},
		},
	)

	return err
}
