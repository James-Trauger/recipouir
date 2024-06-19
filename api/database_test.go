package main

import (
	"context"
	"testing"

	//recapi "github.com/James-Trauger/Recipouir/api"
	"github.com/James-Trauger/Recipouir/model"
	"go.mongodb.org/mongo-driver/bson"
)

func insertUser(uname, pass string) *model.User {
	user := model.NewUser(uname, pass)
	if user == nil {
		return nil
	}
	userCollection.InsertOne(context.Background(), user)
	return user
}

func TestFind(t *testing.T) {
	t.Parallel()
	user := insertUser("test-name", "pass")
	if user == nil {
		t.Fatal("couldn't insert a user")
	}
	res := userCollection.FindOne(context.Background(), bson.M{"_id": user.ID})
	var userAfterInsert model.User
	res.Decode(&userAfterInsert)

	if res.Err() != nil {
		t.Fatal(res.Err())
	} else if !user.Equal(userAfterInsert) {
		t.Fatal("users are not the same")
	}

	// delete the user
	if err := DeleteUser(&userAfterInsert, context.Background()); err != nil {
		t.Fatal("couldn't delete user: " + err.Error())
	}
}
