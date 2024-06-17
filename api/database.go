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
	var err error = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Couldn't load .env file -> " + err.Error())
	}
	mongoDB := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDB))
	if err != nil {
		log.Fatal(err)
	}

	//
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	//
	if client != nil {
		log.Println("MongoDB connected")
	}
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, db string, collectionName string) *mongo.Collection {
	//var collection *mongo.Collection = client.Database(db).Collection(collectionName)
	//return collection
	return client.Database(db).Collection(collectionName)
}
