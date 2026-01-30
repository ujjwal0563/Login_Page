package models

import (
   "time"
)
type User struct {
	Name        string             `bson:"name"`
	Surname     string             `bson:"surname"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	PhoneNumber string             `bson:"phone_number"`
	CreatedAt   time.Time          `bson:"created_at"`
}
