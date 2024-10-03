package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/James-Trauger/Recipouir/utils"
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
	var styles string = "./styles"
	mux := http.NewServeMux()

	// css stylesheet handler
	mux.Handle("/styles/", http.StripPrefix("/styles", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer utils.DrainClose(r.Body)

		w.Header().Set("Content-Type", "text/css")
		// write the file
		fmt.Print(filepath.Join(styles, r.URL.Path))
		http.ServeFile(w, r, filepath.Join(styles, r.URL.Path))
	})))

	// root handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Page(templates.NewNav(routes, "home"),
			[]string{"/styles/topnav.css"}, templates.Welcome()).Render(r.Context(), w)
	})

	// my recipes handler
	mux.HandleFunc("/myrecipes", func(w http.ResponseWriter, r *http.Request) {
		templates.Page(templates.NewNav(routes, "my recipes"),
			[]string{"/styles/topnav.css"}, templates.RecipePage()).Render(r.Context(), w)
	})
	//fmt.Println(http.ListenAndServeTLS(*addr, *cert, *pkey, mux))
	fmt.Println(http.ListenAndServeTLS(*addr, *cert, *pkey, mux))
}
