package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/James-Trauger/Recipouir/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// open the user collection from the db database
var userCollection *mongo.Collection = OpenCollection(Client, DbName, "user")

func Signup(r *http.Request, ctx context.Context) (*model.User, int, error) {

	login, err := model.ExtractLogin(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// should return an error because the user doesn't exist yet
	_, err = Authenticate(login, ctx)

	// user already exists
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, http.StatusUnauthorized, err
	}

	// add the user to the database
	user := model.NewUser(login.Uname, login.Pass)
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		//return nil, http.StatusInternalServerError, errors.New("couldn't insert new user into db")
		return nil, http.StatusInternalServerError, err
	}
	//user.ID = result.InsertedID.(primitive.ObjectID)
	user.ID = result.InsertedID.(string)
	return user, http.StatusOK, nil
}

func Login(r *http.Request, ctx context.Context) (*model.User, int, error) {
	// retrieve username and password from the request
	login, err := model.ExtractLogin(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	user, isAuthenticated := Authenticate(login, ctx)

	if isAuthenticated == nil {
		return user, http.StatusOK, nil
	} else {
		//return nil, http.StatusUnauthorized, errors.New("incorrect username or password")
		return nil, http.StatusUnauthorized, isAuthenticated
	}
}

/* retrieves the target user *
func GetUser(login model.Login) *model.User {

	// give the query a 10 second time limit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filer := bson.M{"_id": login.Uname}
	query := userCollection.FindOne(ctx, filer)

	var user model.User
	query.Decode(&user)

	Authenticate()

	return &user
}*/

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
