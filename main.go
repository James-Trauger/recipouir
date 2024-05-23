package main

import (
	"net/http"

	"github.com/James-Trauger/Recipouir/templates"
	"github.com/a-h/templ"
)

func main() {
	/*var routes map[string]templ.SafeURL = map[string]templ.SafeURL{
		"home":       templ.SafeURL("/home"),
		"my recipes": templ.SafeURL("/my-recipes"),
	}*/
	comp := templates.Index("home")
	//cs := templ.NewCSSHandler()
	http.Handle("/", templ.Handler(comp))
	http.ListenAndServe(":3999", nil)
}
