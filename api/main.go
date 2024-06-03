package main

import (
	"fmt"
	"net/http"
)

const port = "9876"

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api", rootHandler())

	fmt.Println(http.ListenAndServe("127.0.0.1:"+port, mux))
}
