package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	reciauth "github.com/James-Trauger/Recipouir/auth"
	"github.com/James-Trauger/Recipouir/model"
	"github.com/James-Trauger/Recipouir/utils"
)

type tokenUnameKey int

const (
	userKey tokenUnameKey = iota
	loginKey
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
			ErrInternalServer.WriteError(w)
		} else {
			// pass on the the next handler
			handler.ServeHTTP(w, r)
		}
		return
	}

	// method is not supported, key doesn't exist in the map
	w.Header().Add("Allow", utils.AllowedMethods(rm))
	if r.Method != http.MethodOptions {
		NewJsonErr("Method not allowed", http.StatusMethodNotAllowed).WriteError(w)
		return
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
			NewJsonErr(err.Error(), http.StatusBadRequest).WriteError(w)
			return
		}

		// claims of the token
		claims, err := reciauth.ValidToken(token)
		if err != nil {
			NewJsonErr("invalid token", http.StatusBadRequest).WriteError(w)
			return
		}

		// context with the username of who owns the token
		userCtx := context.WithValue(r.Context(), userKey, claims.Username)
		// add the context to the request
		r = r.WithContext(userCtx)

		next.ServeHTTP(w, r)
	})
}

// add the login object to request's context assuming it exists
func loginExtractor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var login model.Login
		// add a value to the request's context
		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			next.ServeHTTP(w, r)
		} else {
			loginCtx := context.WithValue(r.Context(), loginKey, login)
			next.ServeHTTP(w, r.WithContext(loginCtx))
		}
	})
}

/* handles login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the username from the jwt token
	tokenUname := r.Context().Value(userKey)
	tokenUname, ok := tokenUname.(string)
	// cast username to string
	if !ok {
		// username from the token is not a string
		JSONError(w, http.StatusBadRequest, errors.New("inavlid username from jwt token"))
		return
	}

	// copy of the body for
	body, err := r.GetBody()
	if err != nil {
		JSONError(w, http.StatusInternalServerError, errors.New("interal server error"))
		return
	}
	// retrieve credentials from the request
	reqLogin, err := model.ExtractLogin(body)
	if err != nil {
		//TODO bad request, invalid format
		return
	}

	// authorize the user
	if tokenUname != reqLogin.Uname {
		JSONError(w, http.StatusUnauthorized, errors.New("can't login as that user with the token provided"))
		return
	}

}*/
