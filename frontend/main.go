package main

import (
	"flag"
	"fmt"
	"net/http"
	"regexp"
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

	mux.Handle("/assets", mimeMiddleware(http.FileServer(http.Dir("./assets"))))

	for _, r := range routes {
		mux.Handle(r, rootHandler())
	}

	fmt.Println(http.ListenAndServeTLS(*addr, *cert, *pkey, mux))
}

func mimeMiddleware(next http.Handler) http.Handler {
	// extract the file extension from a file
	js, _ := regexp.Compile(".js$")
	css, _ := regexp.Compile(".css$")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path

		next.ServeHTTP(w, r)

		/*parts := strings.Split(p, "/")
		// extract the last segment of the path
		fileName := parts[len(parts)-1]*/
		fmt.Println(p)
		if js.MatchString(p) {
			w.Header().Set("Content-Type", "text/javascript")
		} else if css.MatchString(p) {
			w.Header().Set("Content-Type", "text/css")
		}

	})
}

func rootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		http.ServeFile(w, r, "./views/index.html")
	})
}
