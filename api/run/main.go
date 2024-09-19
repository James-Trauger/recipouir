package main

import (
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

	//cert := os.Getenv("CERT")
	port := os.Getenv("PORT")
	//port := "9872"
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
