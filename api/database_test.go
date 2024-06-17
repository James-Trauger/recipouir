package main_test

import (
	"context"
	"testing"

	recapi "github.com/James-Trauger/Recipouir/api"
	"github.com/James-Trauger/Recipouir/model"
	"go.mongodb.org/mongo-driver/bson"
)

var userCollection = recapi.OpenCollection(recapi.Client, "db", "user")

func insertUser(uname, pass string) *model.User {
	user := model.NewUser(uname, pass)
	userCollection.InsertOne(context.Background(), user)
	return user
}

func TestFind(t *testing.T) {
	user := insertUser("test-name", "pass")

	res := userCollection.FindOne(context.Background(), bson.M{"_id": user.ID})
	var userAfterInsert model.User
	res.Decode(&userAfterInsert)

	if res.Err() != nil {
		t.Fatal(res.Err())
	} else if !user.Equal(userAfterInsert) {
		t.Fatal("users are not the same")
	}
}

/*
func TestFindAll(t *testing.T) {
	userCollection := recapi.OpenCollection(recapi.Client, "db", "user")

	res, err := userCollection.Find(context.Background(), bson.M{})

	if err != nil {
		t.Fatal(err)
	}

	for res.Next(context.Background()) {
		fmt.Println(res.Current.String())
	}
}
*/
