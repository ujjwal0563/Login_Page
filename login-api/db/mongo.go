package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection

func ConnectMongo() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not set")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	UserCollection = client.Database("authdb").Collection("users")
	log.Println("MongoDB connected ")
}
