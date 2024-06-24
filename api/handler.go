package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	reciauth "github.com/James-Trauger/Recipouir/auth"
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
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
			// add the token to the header
			w.Header().Set("content-type", "application/jwt")
			fmt.Fprintln(w, signed)

		}),
	}
}
