package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/James-Trauger/Recipouir/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// open the user collection from the db database
var userCollection *mongo.Collection = OpenCollection(Client, "db", "user")

// returns nil on a successful login: username exists and the password was corrects
func withUsername(name, pass string, ctx context.Context) error {

	filter := bson.M{`"username"`: "\"" + name + "\""}
	result := userCollection.FindOne(ctx, filter)
	var user model.User
	err := result.Decode(&user)
	if err != nil {
		// couldn't decod the user
		//return errors.New("internal server error")
		return err
	}
	return bcrypt.CompareHashAndPassword(*user.Pass, []byte(pass))
}

func Signup(r *http.Request, ctx context.Context) (*model.User, int, error) {

	login, status, err := Login(r, ctx)

	// user already exists
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, status, err
	}

	// add the user to the database
	user := model.NewUser(login.Uname, login.Pass)
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		//return nil, http.StatusInternalServerError, errors.New("couldn't insert new user into db")
		return nil, http.StatusInternalServerError, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, http.StatusOK, nil
}

func Login(r *http.Request, ctx context.Context) (*model.Login, int, error) {
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

	auth := withUsername(login.Uname, login.Pass, ctx)

	if auth == nil {
		return &login, http.StatusOK, nil
	} else {
		//return nil, http.StatusUnauthorized, errors.New("incorrect username or password")
		return &login, http.StatusUnauthorized, auth
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

// returns nil on success
func DeleteUser(target *model.User, ctx context.Context) error {
	bs, err := bson.Marshal(target)
	if err != nil {
		return err
	}

	result, err := userCollection.DeleteOne(ctx, bs)
	if result.DeletedCount != 1 {
		return errors.New("more than one user delete")
	}
	return err
}
