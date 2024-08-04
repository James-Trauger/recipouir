package main

import (
	"log"
	"net/http"

	"github.com/James-Trauger/Recipouir/api"
)

const (
	port = "9876"
	host = "127.0.0.1:"
)

func main() {
	mux := http.NewServeMux()

	// register all routes TODO change recipe model's _id
	mux.Handle("/api/", api.RootHandler())
	mux.Handle(api.LoginPath, api.HandleLogin())
	mux.Handle(api.SignupPath, api.SignupHandler())
	mux.Handle(api.DeleteUserPath, api.DeleteUserHandler())
	mux.Handle(api.AddRecipePath, api.AddRecipeHandler())
	mux.Handle(api.GetRecPath, api.GetRecipeURLHandler())
	mux.Handle(api.GetAllRecPath, api.GetUserRecipesHandler())

	log.Println(http.ListenAndServe(host+port, mux))
}
