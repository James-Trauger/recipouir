package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/James-Trauger/Recipouir/utils"
)

const (
	SignupRoute = "user/signup"
	LoginRoute  = "user/login"
)

func rootHandler() http.Handler {
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
			authErr := utils.Authorize(&r.Header, username)
			if authErr != nil {
				JSONError(w, http.StatusUnauthorized, authErr)
			}
		}),
	}
}
