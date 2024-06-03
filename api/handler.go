package main

import (
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
