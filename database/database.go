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

var Client *mongo.Client = DBinstance()
var DbName = "test"

func DBinstance() *mongo.Client {
	if err := godotenv.Load("../../.env"); err != nil {
		godotenv.Load("../.env")
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
	if client != nil {
		log.Println("MongoDB connected at " + mongoDB)
	}
	return client
}

func OpenCollection(client *mongo.Client, db string, collectionName string) *mongo.Collection {
	return client.Database(db).Collection(collectionName)
}
