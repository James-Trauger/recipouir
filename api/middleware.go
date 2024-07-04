package api

import (
	"errors"
	"io"
	"net/http"

	reciauth "github.com/James-Trauger/Recipouir/auth"
	"github.com/James-Trauger/Recipouir/model"
	"github.com/James-Trauger/Recipouir/utils"
)

type RestMethods map[string]http.Handler

func (rm RestMethods) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// close and drain the request so the tcp connection can be reused
	defer utils.DrainClose(r.Body)

	// call the handler corresponding to the request
	if handler, ok := rm[r.Method]; ok {
		// handler not implemented
		if handler == nil {
			// error on the server side
			JSONError(w, http.StatusInternalServerError, errors.New("Internal server error"))
		} else {
			// pass on the the next handler
			handler.ServeHTTP(w, r)
		}
		return
	}

	// method is not supported, key doesn't exist in the map
	w.Header().Add("Allow", utils.AllowedMethods(rm))
	if r.Method != http.MethodOptions {
		JSONError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
	}
}

// drain and close the request so the tcp session can be reused
func drainAndClose(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
	})
}

func authenticateLogin(next http.Handler, tokenUser string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		model.ExtractLogin(r.Body)
	})
}

func validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract the token from the request
		token, err := reciauth.ParseTokenFromHeader(&r.Header)
		// invalid token
		if err != nil {
			JSONError(w, http.StatusBadRequest, errors.New("couldn't parse token"))
			return
		}

		// claims of the token
		claims, err := reciauth.ValidToken(token)
		if err != nil {
			JSONError(w, http.StatusBadRequest, errors.New("invalid token"))
			return
		}
		authenticateLogin(next, claims.Username)
	})
}
