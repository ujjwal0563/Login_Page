package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"login-api/models"
)

//  Find user by email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := UserCollection.FindOne(
		context.TODO(),
		bson.M{"email": email},
	).Decode(&user)

	return user, err
}

//  Save reset token
func SaveResetToken(
	userID primitive.ObjectID,
	token string,
	expiry time.Time,
) error {

	_, err := UserCollection.UpdateByID(
		context.TODO(),
		userID,
		bson.M{
			"$set": bson.M{
				"reset_token":        token,
				"reset_token_expiry": expiry,
			},
		},
	)

	return err
}