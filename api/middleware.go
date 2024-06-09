package main

import (
	"errors"
	"net/http"

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

/* authentication
func AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract token
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		// tokens are in the form of `Bearer 0x...`
		if len(authHeader) != 2 || authHeader[0] != "Bearer" {
			JSONError(w, http.StatusBadRequest, errors.New("malformed authorization header, expected \"authorization: Bearer [token]\""))
			return // invalid header
		}
		token, err := jwt.ParseWithClaims(authHeader[1], &UserClaims{}, utils.VerifyToken)
		if err != nil {
			JSONError(w, http.StatusInternalServerError, errors.New("couldn't parse token -> "+err.Error()))
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			JSONError(w, http.StatusBadRequest, errors.New("couldn't validate/verify token"))
		}
	})
}
*/
