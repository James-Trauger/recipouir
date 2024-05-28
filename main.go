package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/James-Trauger/Recipouir/templates"
)

var (
	addr = flag.String("listen", "127.0.0.1:8443", "listen address")
	cert = flag.String("cert", "", "certificate")
	pkey = flag.String("key", "", "private key")
)

func main() {
	flag.Parse()
	/*var routes map[string]templ.SafeURL = map[string]templ.SafeURL{
		"home":       templ.SafeURL("/home"),
		"my recipes": templ.SafeURL("/my-recipes"),
	}*/
	mux := http.NewServeMux()

	comp := templates.Index("home")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		comp.Render(r.Context(), w)
	})

	mux.HandleFunc("/", handler)
	mux.HandleFunc("/styles/main.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/topnav.css")
	})
	//cs := templ.NewCSSHandler()
	//http.Handle("/", templ.Handler(handler))
	fmt.Println(http.ListenAndServeTLS(*addr, *cert, *pkey, mux))
}
