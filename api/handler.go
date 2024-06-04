package main

import (
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
