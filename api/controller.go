package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/James-Trauger/Recipouir/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// open the user collection from the db database
var userCollection *mongo.Collection = OpenCollection(Client, "db", "user")

func HashPassword() {

}

func VerifyPassword()

func Signup()

func Login()

func GetUsers()

/*
	/api/user?userid=...&userType=...

only called after the user is authenticated
*/
func GetUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		uid := query.Get("userid")
		if uid == "" {
			JSONError(w, http.StatusBadRequest, errors.New("no userid provided"))
			return
		}

		//TODO authorize user

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user model.User
		err := userCollection.FindOne(ctx, bson.M{"userid": uid}).Decode(&user)
		defer cancel()
		if err != nil {
			JSONError(w, http.StatusInternalServerError, errors.New("user not found"))
		}
	})
}
