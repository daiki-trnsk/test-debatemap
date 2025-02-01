package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not set")
	}
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatalf("Failed to ping MongoDB: %v", err)
    }
	fmt.Println("Successfully Connected to the mongodb")
	return client
}

var Client *mongo.Client = DBSet()

func DebateMapData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var debatemapcollection *mongo.Collection = client.Database("testDebateMap").Collection(CollectionName)
	return debatemapcollection
}
