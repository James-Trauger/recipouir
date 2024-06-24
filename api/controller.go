package api

import (
	"context"
	"errors"
	"net/http"

	reciauth "github.com/James-Trauger/Recipouir/auth"
	db "github.com/James-Trauger/Recipouir/database"
	"github.com/James-Trauger/Recipouir/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// open the user collection from the db database
var (
	userCollection   *mongo.Collection = db.OpenCollection(db.Client, db.DbName, "user")
	recipeCollection *mongo.Collection = db.OpenCollection(db.Client, db.DbName, "recipe")
)

func Signup(r *http.Request, ctx context.Context) (*model.User, int, error) {

	login, err := model.ExtractLogin(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// should return an error because the user doesn't exist yet
	_, err = reciauth.Authenticate(login, ctx)

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

	user, isAuthenticated := reciauth.Authenticate(login, ctx)

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

/*
reutrns the recipe from a specific user. Recipes are public so no authentication
nor authorization is needed besides having a valid jwt token
*/
func GetRecipe(user, name string, ctx context.Context) (*model.Recipe, error) {

	//TODO validate token

	// filter based on the name of a recipe the specified user created
	filter := bson.M{"user": user, "name": name}
	var recipe model.Recipe
	result := recipeCollection.FindOne(context.TODO(), filter)
	err := result.Decode(&recipe)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

// inserts a single recipe assuming the user is authorized and returns nil on succes
func InsertRecipe(rec model.Recipe, user string, ctx context.Context) error {
	_, err := recipeCollection.InsertOne(ctx, rec)
	return err
}

// insert multiple recipes assuming the user is authroized, returns nil on success
func InsertManyRecipe(recipes *[]model.Recipe, user string, ctx context.Context) error {
	cast := []any{*recipes} // TODO test this, prolly wont work
	_, err := recipeCollection.InsertMany(ctx, cast)
	return err
}
