package utils

import (
	"io"
	"net/http"
)

func DrainClose(r io.ReadCloser) error {
	io.Copy(io.Discard, r)
	return r.Close()
}

/* middleware for linking stylesheets to the page */
func MiddleCSS(next http.Handler, styles []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
