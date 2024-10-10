package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/James-Trauger/Recipouir/web/templates"
	"github.com/a-h/templ"
)

var (
	addr = flag.String("listen", "127.0.0.1:8443", "listen address")
	cert = flag.String("cert", "", "certificate")
	pkey = flag.String("key", "", "private key")
)

func main() {
	flag.Parse()
	var routes map[string]templ.SafeURL = map[string]templ.SafeURL{
		"home":       templ.SafeURL("/"),
		"my recipes": templ.SafeURL("/myrecipes"),
	}
	//var styles string = "./styles"
	mux := http.NewServeMux()

	// root handler
	mux.Handle(string(routes[rootName]), rootHandler())

	// my recipes handler
	mux.HandleFunc("/myrecipes", func(w http.ResponseWriter, r *http.Request) {
		templates.Page(templates.NewNav(routes, "my recipes"),
			[]string{}, //[]string{"/styles/topnav.css"},
			templates.RecipePage()).Render(r.Context(), w)
	})

	fmt.Println(http.ListenAndServeTLS(*addr, *cert, *pkey, mux))
}
