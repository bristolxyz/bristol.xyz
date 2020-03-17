package clients

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDatabase is the MongoDB database we are using.
var MongoDatabase *mongo.Database

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
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		return err
	}
	MongoDatabase = client.Database(DB)
	return client.Ping(context.TODO(), readpref.Primary())
}
