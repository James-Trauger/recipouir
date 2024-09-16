package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func ClientTest(t *testing.T) {
	api := "localhost:9872/api/"
	resp, err := http.Get(api)
	if err != nil {
		fmt.Println(resp.Status)
		t.Fatal(err)
	}
	io.Copy(os.Stdout, resp.Body)
}
