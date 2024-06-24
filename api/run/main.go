package api

import (
	"fmt"
	"net/http"

	"github.com/James-Trauger/Recipouir/api"
)

const port = "9876"

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api", api.RootHandler())

	fmt.Println(http.ListenAndServe("127.0.0.1:"+port, mux))
}
