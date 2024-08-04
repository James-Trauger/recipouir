package api

import (
	"context"
	"errors"

	reciauth "github.com/James-Trauger/Recipouir/auth"
	db "github.com/James-Trauger/Recipouir/database"
	"github.com/James-Trauger/Recipouir/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// open the user collection from the db database
var (
	userCollection      *mongo.Collection = db.OpenCollection(db.Client, db.DbName, "user")
	recipeCollection    *mongo.Collection = db.OpenCollection(db.Client, db.DbName, "recipe")
	ErrUserAlreadyExits                   = errors.New("user already exists")
	ErrDbInsertionError                   = errors.New("couldn't insert the passed data")
	ErrRecipeExists                       = errors.New("recipe already exists")
)

// adds a user to the database based on the provided credentials assuming the
// username does not exist yet.
func Signup(login *model.Login, ctx context.Context) (*model.User, error) {
	// should return an error because the user doesn't exist yet
	_, err := reciauth.Authenticate(login, ctx)

	// user already exists
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserAlreadyExits
	}

	// add the user to the database
	user := model.NewUser(login.Uname, login.Pass)
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, ErrDbInsertionError
	}
	//user.ID = result.InsertedID.(primitive.ObjectID)
	user.ID = result.InsertedID.(string)
	return user, nil
}

// validate a user's credentials and returns a valid jwt token
func Login(creds *model.Login, ctx context.Context) (*model.User, error) {
	// validate the credentials
	user, isAuthenticated := reciauth.Authenticate(creds, ctx)

	if isAuthenticated == nil {
		return user, nil
	} else {
		//return nil, http.StatusUnauthorized, errors.New("incorrect username or password")
		return nil, isAuthenticated
	}
}

// returns nil on success
// delete a user and their recipe data from the database
func DeleteUser(target *model.User, ctx context.Context) error {
	bs, err := bson.Marshal(target)
	if err != nil {
		return err
	}

	// delete user
	result, err := userCollection.DeleteOne(ctx, bs)
	if result.DeletedCount != 1 {
		return errors.New("couldn't delete user")
	}
	return err
}

// delete all of a user's recipes. used when a user wants to delete their account
func DeleteAllRecipes(user string, ctx context.Context) error {
	// filter on which user created the recipe
	filter := bson.M{"createdby": user}
	_, err := recipeCollection.DeleteMany(ctx, filter)
	return err
}

// reutrns the recipe from a specific user. Recipes are public so no authentication
// nor authorization is needed besides having a valid jwt token
func GetRecipe(user, name string, ctx context.Context) (*model.Recipe, error) {
	// filter based on the name of a recipe the specified user created
	filter := bson.M{"name": name, "createdby": user}
	var recipe model.Recipe
	result := recipeCollection.FindOne(ctx, filter)
	err := result.Decode(&recipe)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

// returns all the recipes that a particular user has
func GetAllRecipes(user string, ctx context.Context) (*[]model.Recipe, error) {

	// filter only on username
	filter := bson.M{"createdby": user}
	result, err := recipeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var recipes []model.Recipe
	if err := result.All(ctx, &recipes); err != nil {
		return nil, err
	}
	return &recipes, nil
}

// inserts a single recipe assuming the user is authorized and returns nil on succes
func InsertRecipe(rec model.Recipe, user string, ctx context.Context) error {
	// see if the recipe already exists
	if res, _ := GetRecipe(rec.CreatedBy, rec.Name, ctx); res != nil {
		return ErrRecipeExists
	}
	_, err := recipeCollection.InsertOne(ctx, rec)
	return err
}

// insert multiple recipes assuming the user is authroized, returns nil on success
func InsertManyRecipe(recipes *[]model.Recipe, user string, ctx context.Context) error {
	cast := make([]any, len(*recipes))
	// cast each element to type any
	for i := range *recipes {
		cast[i] = any((*recipes)[i])
	}

	_, err := recipeCollection.InsertMany(ctx, cast)
	return err
}

// delete a recipe from the database based on its name and user who created it
func DeleteRecipe(name, user string, ctx context.Context) error {
	filter := bson.M{"name": name, "createdby": user}
	_, err := recipeCollection.DeleteOne(ctx, filter)
	return err
}
