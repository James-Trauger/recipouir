package main

import (
	"fmt"
	"net/http"

	"github.com/James-Trauger/Recipouir/api"
)

const (
	port       = "9876"
	rootPath   = "/api"
	signupPath = "/signup"
	loginPath  = "/login"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api", api.RootHandler())
	//mux.Handle("/api/user/", api.RootHandler())

	fmt.Println(http.ListenAndServe("127.0.0.1:"+port, mux))
}
