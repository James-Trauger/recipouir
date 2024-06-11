package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// returns nil on a successful login: username exists and the password was corrects
func withUsername(name, pass string) error {

	filter := bson.M{`"username"`: "\"" + name + "\""}
	result := userCollection.FindOne(context.TODO(), filter)
	var user model.User
	err := result.Decode(&user)
	if err != nil {
		// internal server error
		return errors.New("internal server error")
	}
	return bcrypt.CompareHashAndPassword(*user.Pass, []byte(pass))
}

func Signup(r *http.Request) (*model.Login, int, error) {

	login, status, err := Login(r)

	// user already exists
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, status, err
	}

	// add the user to the database
	user := model.NewUser(*login.Uname, *login.Pass)
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("couldn't insert new user into db")
	}
	return login, http.StatusOK, nil
}

func Login(r *http.Request) (*model.Login, int, error) {
	// try logging the user in to see if they already have an account
	var login model.Login
	buf := bytes.NewBuffer(nil)
	n, err := io.Copy(buf, r.Body)
	if err != nil {
		// body not read
		return nil, http.StatusInternalServerError, err
	} else if n == 0 {
		// empty body
		return nil, http.StatusBadRequest, errors.New("no request body provided")
	}

	// decode the json request
	if err = json.Unmarshal(buf.Bytes(), &login); err != nil {
		return nil, http.StatusBadRequest, err
	}

	auth := withUsername(*login.Uname, *login.Pass)

	if auth == nil {
		return &login, http.StatusOK, nil
	} else {
		return nil, http.StatusUnauthorized, errors.New("incorrect username or password")
	}
}

func HandleLogin() http.Handler {
	return RestMethods{
		http.MethodPost: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			login, status, err := Login(r)

			if err != nil {
				JSONError(w, status, err)
				return
			}

			// valid credentials
			// return a jwt token using RSA, expires a day from now
			signed, err := utils.NewToken(*login.Uname)
			if err != nil {
				JSONError(w, http.StatusInternalServerError, errors.New("couldn't create jwt token"))
			}
			// add the token to the header
			w.Header().Set("content-type", "application/jwt")
			fmt.Fprintln(w, signed)

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
