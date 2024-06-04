package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	var err error = godotenv.Load(".env") // TODO either .env or the path to the .env
	if err != nil {
		log.Fatal("Couldn't load .env file -> " + err.Error())
	}
	MongoDB := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDB))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB connected")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, db string, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(db).Collection(collectionName)
	return collection
}
