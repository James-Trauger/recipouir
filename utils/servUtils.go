package utils

import (
	"io"
	"net/http"
	"sort"
	"strings"
)

// http requests map to a handler
type Methods map[string]http.Handler

func (h Methods) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// close and drain the request so the tcp connection can be reused
	defer DrainClose(r.Body)

	// call the handler corresponding to the request
	if handler, ok := h[r.Method]; ok {
		// handler not implemented
		if handler == nil {
			// error on the server side
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			// pass on the the next handler
			handler.ServeHTTP(w, r)
		}
		return
	}

	// method is not supported, key doesn't exist in the map
	w.Header().Add("Allow", AllowedMethods(h))
	if r.Method != http.MethodOptions {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func AllowedMethods(meth map[string]http.Handler) string {
	// element for each key
	a := make([]string, 0, len(meth))

	// add all the methods
	for k := range meth {
		a = append(a, k)
	}
	sort.Strings(a)
	return strings.Join(a, ", ")
}

func DrainClose(r io.ReadCloser) error {
	io.Copy(io.Discard, r)
	return r.Close()
}

/* middleware for linking stylesheets to the page */
func MiddleCSS(next http.Handler, styles []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
