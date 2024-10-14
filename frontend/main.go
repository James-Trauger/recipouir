package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	addr = flag.String("listen", "127.0.0.1:8443", "listen address")
	cert = flag.String("cert", "", "certificate")
	pkey = flag.String("key", "", "private key")
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()

	routes := []string{
		"/", // home page
		"/about",
	}

	// serve assets built with react
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	for _, r := range routes {
		mux.Handle(r, rootHandler())
	}

	fmt.Println(http.ListenAndServeTLS(*addr, *cert, *pkey, mux))
}

// serves index.tmpl
func rootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		http.ServeFile(w, r, "./views/index.html")
	})
}
