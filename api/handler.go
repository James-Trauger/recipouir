package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	reciauth "github.com/James-Trauger/Recipouir/auth"
	"github.com/James-Trauger/Recipouir/model"
	"github.com/James-Trauger/Recipouir/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	SignupRoute      = "user/signup"
	LoginRoute       = "user/login"
	userPatternURL   = "username"
	recipePatternURL = "recipe"
)

var (
	ErrInternalServer     = errors.New("internal server error")
	ErrNotFound           = errors.New("route does not exist")
	ErrNoRecipe           = errors.New("recipe not found")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

func RootHandler() http.Handler {
	return utils.Methods{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("test"))
		}),
	}
}

func JSONError(w http.ResponseWriter, code int, err error) (int, error) {
	w.Header().Set("content-type", "application/json")
	body := []byte(fmt.Sprintf("{\"error\":\"%s\"}\n", err))
	w.WriteHeader(code)
	return w.Write(body)
}

// returns the user within the query
// /api/user?username=...
func userHandler() http.Handler {
	return RestMethods{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := r.URL.Query().Get("username")
			if username == "" {
				JSONError(w, http.StatusBadRequest, errors.New("no username query"))
				return
			}
			authErr := reciauth.Authorize(&r.Header, username)
			if authErr != nil {
				JSONError(w, http.StatusUnauthorized, authErr)
			}
		}),
	}
}

type tokenResponse struct {
	Token string `json:"token"`
}

// used to get jwt token
func HandleLogin() http.Handler {
	return RestMethods{
		http.MethodPost: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()

			creds, err := model.ExtractLogin(r.Body)
			if err != nil {
				// TODO
				JSONError(w, http.StatusBadRequest, err)
				return
			}

			user, err := Login(creds, ctx)

			if err != nil {
				// TODO check error type
				JSONError(w, http.StatusBadRequest, ErrInvalidCredentials)
				return
			}

			// valid credentials
			// return a jwt token using RSA, expires a day from now
			signed, err := reciauth.NewToken(user.Username)
			if err != nil {
				JSONError(w, http.StatusInternalServerError, errors.New("couldn't create jwt token"))
			}
			// add the token
			w.Header().Set("content-type", "application/json")

			tokResp := tokenResponse{Token: signed}
			json.NewEncoder(w).Encode(tokResp)
			//fmt.Fprintf(w, "{token: %s}", signed)

		}),
	}
}

func SignupHandler() http.Handler {
	return RestMethods{
		http.MethodPost: loginExtractor(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			login, ok := ctx.Value(loginKey).(model.Login)
			if !ok {
				JSONError(w, http.StatusBadRequest, errors.New("no username or password provided"))
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
				JSONError(w, status, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			// TODO Dsuccessful signup message?
		})),
	}
}

// delete a user and all their recipes
func DeleteUserHandler() http.Handler {
	return RestMethods{
		http.MethodPost: validateToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()

			// extract user
			creds, err := model.ExtractLogin(r.Body)
			if err != nil {
				// TODO
				JSONError(w, http.StatusBadRequest, err)
				return
			}

			// validate credentials
			user, err := Login(creds, ctx)
			if err != nil {
				JSONError(w, http.StatusBadRequest, errors.New("invalid credentials"))
				return
			}

			// delete the recipe data in db
			if DeleteAllRecipes(user.Username, ctx) != nil {
				JSONError(w, http.StatusInternalServerError, ErrInternalServer)
				return
			}
			// delete user in db
			if DeleteUser(user, ctx) != nil {
				JSONError(w, http.StatusInternalServerError, ErrInternalServer)
				return
			}

			w.WriteHeader(http.StatusOK)
		})),
	}
}

func AddRecipeHandler() http.Handler {
	return RestMethods{
		http.MethodPost: validateToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract the recipe
			var recipe model.Recipe
			err := json.NewDecoder(r.Body).Decode(&recipe)
			if err != nil {
				JSONError(w, http.StatusBadRequest, errors.New("couldn't decode the recipe"))
				return
			}

			// get the username from the token
			userToken := r.Context().Value(userKey)
			if userToken == nil {
				JSONError(w, http.StatusInternalServerError, errors.New("couldn't find user associated with the token"))
				return
			}

			// token must match the user adding the recipe
			if userToken != recipe.CreatedBy {
				JSONError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
				return
			}

			// give the database 5 seconds to insert the recipe
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			// add the recipe
			if err = InsertRecipe(recipe, recipe.CreatedBy, ctx); err != nil {
				JSONError(w, http.StatusInternalServerError, err)
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
			user, recipe := r.PathValue(userPatternURL), r.PathValue(recipePatternURL)
			if user == "" || recipe == "" {
				JSONError(w, http.StatusNotFound, ErrNotFound)
			}

			rec, err := GetRecipe(user, recipe, r.Context())
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					JSONError(w, http.StatusNoContent, ErrNoRecipe)
				}
				// TODO get recipe error, check error type
				JSONError(w, http.StatusInternalServerError, ErrInternalServer)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "application/json")
			// return the recipe
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(&rec)
			if err != nil {
				JSONError(w, http.StatusInternalServerError, ErrInternalServer)
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
			user := r.PathValue(userPatternURL)
			if user == "" {
				JSONError(w, http.StatusNotFound, ErrNotFound)
				return
			}

			recs, err := GetAllRecipes(user, r.Context())
			if err != nil {
				if len(*recs) == 0 {
					JSONError(w, http.StatusNoContent, ErrNoRecipe)
				}
				// TODO check error type
				JSONError(w, http.StatusInternalServerError, ErrInternalServer)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "application/json")
			err = json.NewEncoder(w).Encode(recs)

			if err != nil {
				JSONError(w, http.StatusInternalServerError, ErrInternalServer)
				return
			}
		}),
	}
}
