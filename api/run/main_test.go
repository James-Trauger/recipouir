package main

/*
import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	t.Parallel()
	done := make(chan struct{})

	mux := http.NewServeMux()
	mux.Handle("/api/user/{username}/{recipe}/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		fmt.Println(r.URL.Path)
		fmt.Println(r.PathValue("username"))
		fmt.Println(r.PathValue("recipe"))
		done <- struct{}{}
	}))

	go func() { fmt.Println(http.ListenAndServe("127.0.0.1:9999", mux)) }()

	time.Sleep(time.Second)
	go fmt.Println(http.Get("http://127.0.0.1:9999/api/user/trau/brownies"))
	//req := httptest.NewRequest(http.MethodGet, "https://recipouir.com/user/trau/my-recipe", nil)
	//paths := strings.Split(req.URL.Path, "/")
	//fmt.Println(paths)
	<-done
}
*/
