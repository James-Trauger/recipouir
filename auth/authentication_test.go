package reciauth

import (
	"context"
	"testing"
	"time"

	db "github.com/James-Trauger/Recipouir/database"
	"github.com/James-Trauger/Recipouir/model"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAuthenticate(t *testing.T) {
	login := model.Login{
		Uname: "ned",
		Pass:  "honor",
	}
	user := model.NewUser(login.Uname, login.Pass)
	userCollection := db.OpenCollection(db.Client, "test", "user")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// insert the user
	if _, err := userCollection.InsertOne(ctx, user); err != nil {
		t.Fatal(err)
	}

	// authenticate as the inserted user
	foundUser, err := Authenticate(&login, ctx)
	if err != nil {
		t.Fatal(err)
	}

	// delete the inserted user
	userCollection.DeleteOne(ctx, bson.M{"_id": user.Username})

	// compare the found user and the original
	if !foundUser.Equal(*user) {
		t.Fatal("original user is not the same as the authentiated user")
	}

	// try to authenticate with incorrect credentials
	incorrectLogin := model.Login{
		Uname: "ned",
		Pass:  "honour",
	}

	foundUser, err = Authenticate(&incorrectLogin, ctx)
	if err == nil || foundUser != nil {
		t.Fatal("incorrect credentials successfully authenticated")
	}

	t.Logf("SUCCESS: %s", err)
}
