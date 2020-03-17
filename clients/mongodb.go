package clients

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoClient is the MongoDB client we are using.
var MongoClient *mongo.Client

// CreateMongoClient is used to create the MongoDB client.
func CreateMongoClient() error {
	URI := os.Getenv("MONGODB_URI")
	if URI == "" {
		URI = "mongodb://localhost:27017"
	}
	DB := os.Getenv("MONGODB_DATABASE")
	if DB == "" {
		DB = "bristolxyz"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	MongoClient = client
	return client.Ping(context.TODO(), readpref.Primary())
}
