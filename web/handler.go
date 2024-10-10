package main

import (
	"net/http"

	"github.com/James-Trauger/Recipouir/utils"
	"github.com/James-Trauger/Recipouir/web/templates"
	"github.com/a-h/templ"
)

// paths
const (
	rootName   = "home"
	createName = "create"
)

// name of the navbar page mapped to its route
var routes map[string]templ.SafeURL = map[string]templ.SafeURL{
	rootName:   "/",
	createName: "/create",
}

type WebMethods map[string]http.Handler

func (wm WebMethods) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// close and drain the request so the tcp connection can be reused
	defer utils.DrainClose(r.Body)

	// call the handler corresponding to the request
	if handler, ok := wm[r.Method]; ok {
		// handler not implemented
		if handler == nil {
			// error on the server side
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			// pass on the the next handler
			handler.ServeHTTP(w, r)
		}
		return
	}

	// method is not supported, key doesn't exist in the map
	w.Header().Add("Allow", utils.AllowedMethods(wm))
	if r.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func rootHandler() http.Handler {
	return WebMethods{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			templates.Page(templates.NewNav(routes, rootName),
				[]string{}, //[]string{"/styles/topnav.css"},
				templates.Welcome()).Render(r.Context(), w)
		}),
	}
}
