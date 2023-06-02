package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This function connects to a MongoDB database using a provided connection URI and returns a client object.
func ConnectToDatabase() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error while loading ENV %v", err)
	}
	var url string

	if url = os.Getenv("MONGO_CONNECTION_URI"); url == "" {
		log.Fatal("connection url is empty")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(url))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected successfully")

	return client
}

// `var Client *mongo.Client = ConnectToDatabase()` is initializing a global variable `Client` of type
// `*mongo.Client` with the value returned by the `ConnectToDatabase()` function. This variable can be
// accessed and used throughout the package to interact with the MongoDB database.
var Client *mongo.Client = ConnectToDatabase()

// The function returns a MongoDB collection given a client and collection name.
func OpenCollection(client mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("http-crud").Collection(collectionName)
}
