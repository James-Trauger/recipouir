package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	reciauth "github.com/James-Trauger/Recipouir/auth"
	"github.com/James-Trauger/Recipouir/model"
	"github.com/James-Trauger/Recipouir/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	SignupPath     = "/api/signup"
	LoginPath      = "/api/login"
	DeleteUserPath = "/api/user/remove"
	AddRecipePath  = "/api/recipe/add"

	userPatternUrlKey   = "username"
	recipePatternUrlKey = "recipe"
	GetRecPath          = "/api/user/{" + userPatternUrlKey + "}/{" + recipePatternUrlKey + "}"
	GetAllRecPath       = "/api/user/{" + userPatternUrlKey + "}"
)

func RootHandler() http.Handler {
	return utils.Methods{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("test"))
		}),
	}
}

// used to get jwt token
func HandleLogin() http.Handler {
	return RestMethods{
		http.MethodPost: authenticateLogin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			user, ok := ctx.Value(userKey).(model.User)
			if !ok {
				ErrInternalServer.WriteError(w)
				return
			}

			// valid credentials
			// return a jwt token using RSA, expires a day from now
			signed, err := reciauth.NewToken(user.Username)
			if err != nil {
				//TODO
				ErrInternalServer.WriteError(w)
				return
			}
			// add the token
			w.Header().Set("content-type", "application/json")

			tokResp := TokenResponse{Token: signed}
			json.NewEncoder(w).Encode(tokResp)
			//fmt.Fprintf(w, "{token: %s}", signed)

		})),
	}
}

func SignupHandler() http.Handler {
	return RestMethods{
		http.MethodPost: loginExtractor(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			login, ok := ctx.Value(loginKey).(model.Login)
			if !ok {
				NewJsonErr("no username or password provided", http.StatusBadRequest).WriteError(w)
				return
			}

			// add timeout to the context
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			_, err := Signup(&login, ctx)
			if err != nil {
				var status int = http.StatusBadRequest
				if errors.Is(err, ErrUserAlreadyExits) {
					status = http.StatusConflict
				} else {
					status = http.StatusInternalServerError
				}
				NewJsonErr(err.Error(), status).WriteError(w)
				return
			}

			w.WriteHeader(http.StatusOK)
			// TODO successful signup message?
		})),
	}
}

// delete a user and all their recipes
func DeleteUserHandler() http.Handler {
	return RestMethods{
		http.MethodPost: validateToken(authenticateLogin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			user, ok := ctx.Value(userKey).(model.User)
			if !ok {
				ErrInternalServer.WriteError(w)
				return
			}

			// delete the recipe data in db
			if DeleteAllRecipes(user.Username, ctx) != nil {
				ErrInternalServer.WriteError(w)
				return
			}
			// delete user in db
			if DeleteUser(&user, ctx) != nil {
				ErrInternalServer.WriteError(w)
				return
			}

			w.WriteHeader(http.StatusOK)
		}))),
	}
}

func AddRecipeHandler() http.Handler {
	return RestMethods{
		http.MethodPost: validateToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract the recipe
			var recipe model.Recipe
			err := json.NewDecoder(r.Body).Decode(&recipe)
			if err != nil {
				NewJsonErr(err.Error(), http.StatusBadRequest).WriteError(w)
				return
			}

			// get the username from the token
			userToken := r.Context().Value(userKey)
			if userToken == nil {
				NewJsonErr("couldn't find user associated with the token", http.StatusInternalServerError).WriteError(w)
				return
			}

			// token must match the user adding the recipe
			if userToken != recipe.CreatedBy {
				ErrUnauthorized.WriteError(w)
				return
			}

			// give the database 5 seconds to insert the recipe
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			// add the recipe
			if err = InsertRecipe(recipe, recipe.CreatedBy, ctx); err != nil {
				ErrInternalServer.WriteError(w)
				return
			}

			w.WriteHeader(http.StatusOK)
		})),
	}
}

// /api/user/{username}/{recipe}
// return the recipe of the user at the url
func GetRecipeURLHandler() http.Handler {
	return RestMethods{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, recipe := r.PathValue(userPatternUrlKey), r.PathValue(recipePatternUrlKey)
			if user == "" || recipe == "" {
				ErrNotFound.WriteError(w)
			}

			rec, err := GetRecipe(user, recipe, r.Context())
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					ErrNoRecipe.WriteError(w)
				}
				// TODO get recipe error, check error type
				ErrInternalServer.WriteError(w)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "application/json")
			// return the recipe
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(&rec)
			if err != nil {
				ErrInternalServer.WriteError(w)
				return
			}
		}),
	}
}

// returns all of a user's recipes
// /api/user/username
func GetUserRecipesHandler() http.Handler {
	return RestMethods{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// retrieve the username from the url
			user := r.PathValue(userPatternUrlKey)
			if user == "" {
				ErrNotFound.WriteError(w)
				return
			}

			recs, err := GetAllRecipes(user, r.Context())
			if err != nil {
				if len(*recs) == 0 {
					ErrNoRecipe.WriteError(w)
				}
				// TODO check error type
				ErrInternalServer.WriteError(w)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "application/json")
			err = json.NewEncoder(w).Encode(recs)

			if err != nil {
				ErrInternalServer.WriteError(w)
				return
			}
		}),
	}
}
