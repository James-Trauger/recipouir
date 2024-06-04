package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/James-Trauger/Recipouir/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// open the user collection from the db database
var userCollection *mongo.Collection = OpenCollection(Client, "db", "user")

func HashPassword() {

}

func VerifyPassword()

// uType == "email" OR uType == "username"
func withUname(name string, uType string, pass string) (bool, error) {
	if uType != "email" && uType != "username" {
		return false, errors.New("invalid credential type, must be email or username")
	}

	filter := bson.D{{"\"" + uType + "\"", "\"" + name + "\""}}
	result := OpenCollection(Client, "db", "user").FindOne(context.Background(), filter)
	var user model.User
	err := result.Decode(&user)
	if err != nil {
		// internal server error
		return false, errors.New("internal server error")
	}
	err = bcrypt.CompareHashAndPassword(*user.Pass, []byte(pass))
	return err == nil, err
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

			var isAuthenticated bool
			var authError error
			// username provided
			if login.Uname != nil && login.Email == nil {
				isAuthenticated, authError = withUname(*login.Uname, "username", *login.Pass)
			}

			// email provided
			if login.Email != nil && login.Uname == nil {
				isAuthenticated, authError = withUname(*login.Uname, "email", *login.Pass)
			}

			// valid credentials
			if isAuthenticated {
				// return a jwt token

			} else {
				JSONError(w, http.StatusUnauthorized, authError)
			}
		}),
	}
}

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
