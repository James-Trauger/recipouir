package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/James-Trauger/Recipouir/api"
	"github.com/joho/godotenv"
)

const (
	host = "0.0.0.0:"
)

func main() {
	godotenv.Load("../../.env")
	fmt.Println("mongo db is " + os.Getenv("MONGODB_URL"))
	port := os.Getenv("API_PORT")
	mux := http.NewServeMux()

	mux.Handle("/api/", api.RootHandler())
	mux.Handle(api.LoginPath, api.HandleLogin())
	mux.Handle(api.SignupPath, api.SignupHandler())
	mux.Handle(api.DeleteUserPath, api.DeleteUserHandler())
	mux.Handle(api.AddRecipePath, api.AddRecipeHandler())
	mux.Handle(api.GetRecPath, api.GetRecipeURLHandler())
	mux.Handle(api.GetAllRecPath, api.GetUserRecipesHandler())
	log.Println("listening on ", host, port)
	log.Println(http.ListenAndServe(host+port, mux))
}
