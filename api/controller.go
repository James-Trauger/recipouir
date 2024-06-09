package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/James-Trauger/Recipouir/model"
	"github.com/James-Trauger/Recipouir/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// open the user collection from the db database
var userCollection *mongo.Collection = OpenCollection(Client, "db", "user")

// uType == "email" OR uType == "username"
func withUsername(name, pass string) error {

	filter := bson.M{`"username"`: "\"" + name + "\""}
	result := userCollection.FindOne(context.Background(), filter)
	var user model.User
	err := result.Decode(&user)
	if err != nil {
		// internal server error
		return errors.New("internal server error")
	}
	return bcrypt.CompareHashAndPassword(*user.Pass, []byte(pass))
}

func Signup()

func Login() http.Handler {
	return RestMethods{
		http.MethodPost: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// retrieve the username and password
			var login model.Login
			var buf []byte
			r.Body.Read(buf)
			json.Unmarshal(buf, &login)

			authenticated := withUsername(*login.Uname, *login.Pass)

			// valid credentials
			if authenticated != nil {
				// return a jwt token using RSA, expires a day from now
				signed, err := utils.NewToken(*login.Uname)
				if err != nil {
					JSONError(w, http.StatusInternalServerError, errors.New("couldn't create jwt token"))
				}
				// add the token to the header
				w.WriteHeader(http.StatusOK)
				w.Header().Set("content-type", "application/jwt")
				fmt.Fprintln(w, signed)
			} else {
				// incorrect credentials
				JSONError(w, http.StatusUnauthorized, authenticated)
			}
		}),
	}
}

/* retrieves the target user, they must already be authorized */
func GetUser(target string) *model.User {
	if userCollection == nil {
		return nil
	}
	// give the query a 10 second time limit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filer := bson.M{"username": url.QueryEscape(target)}
	query := userCollection.FindOne(ctx, filer)

	var user model.User
	query.Decode(&user)

	if user.Username != target {
		return nil
	}

	return &user
}

/*
	/api/user?userid=...&userType=...

only called after the user is authenticated
*
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
}*/
