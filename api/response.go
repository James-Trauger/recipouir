package api

import (
	"encoding/json"
	"net/http"
)

var (
	ErrInternalServer = JsonError{
		Msg:  "internal server error",
		code: http.StatusInternalServerError,
	}
	ErrNotFound = JsonError{
		Msg:  "route does not exist",
		code: http.StatusNotFound,
	}
	ErrNoRecipe = JsonError{
		Msg:  "recipe not found",
		code: http.StatusNoContent,
	}
	ErrInvalidCredentials = JsonError{
		Msg:  "invalid username or password",
		code: http.StatusBadRequest,
	}
	ErrUnauthorized = JsonError{
		Msg:  "unauthorized",
		code: http.StatusUnauthorized,
	}
)

type TokenResponse struct {
	Token string `json:"token"`
}

type JsonError struct {
	Msg  string `json:"error"`
	code int
}

func NewJsonErr(msg string, status int) JsonError {
	return JsonError{
		Msg:  msg,
		code: status,
	}
}

func (e *JsonError) Error() string {
	return e.Msg
}

func (e JsonError) WriteError(w http.ResponseWriter) error {
	w.WriteHeader(e.code)
	return json.NewEncoder(w).Encode(e)
}

func (e *JsonError) GetCode() int {
	return e.code
}
