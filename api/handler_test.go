package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/James-Trauger/Recipouir/utils"
)

func TestDefaultHandle(t *testing.T) {
	serv := httptest.NewServer(rootHandler())

	req := httptest.NewRequest(http.MethodGet, serv.URL, nil)
	w := httptest.NewRecorder()
	rootHandler().ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	utils.DrainClose(resp.Body)

	// test method not found
	req = httptest.NewRequest(http.MethodPut, serv.URL, nil)
	w = httptest.NewRecorder()
	rootHandler().ServeHTTP(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		log.Fatal("incorrect response code for an unsupported method")
	}
	utils.DrainClose(resp.Body)

	// test Options method
	req = httptest.NewRequest(http.MethodOptions, serv.URL, nil)
	w = httptest.NewRecorder()
	rootHandler().ServeHTTP(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("options method not supported")
	}
	fmt.Println(resp.Header.Get("Allow"))
	utils.DrainClose(resp.Body)
}
