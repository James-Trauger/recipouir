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
)

const (
	SignupRoute = "user/signup"
	LoginRoute  = "user/login"
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
	body := []byte(fmt.Sprintf("{\"error\":%s}", err))
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

func HandleLogin() http.Handler {
	return RestMethods{
		http.MethodPost: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()

			user, status, err := Login(r, ctx)

			if err != nil {
				JSONError(w, status, err)
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
			fmt.Fprintf(w, "{token: %s}", signed)

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
					status = http.StatusUnprocessableEntity
				} else {
					status = http.StatusInternalServerError
				}
				JSONError(w, status, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			// successful signup message?
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
		})),
	}
}

// return the recipe of the user at the url
func GetRecipeHandler() http.Handler {
	return RestMethods{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		}),
	}
}
