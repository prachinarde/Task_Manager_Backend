package config

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB now returns an error
func ConnectDB() error {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("❌ MongoDB Connection Failed:", err)
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("❌ MongoDB Ping Failed:", err)
		return err
	}

	DB = client.Database("taskManager")
	log.Println("✅ Connected to MongoDB")
	return nil
}
