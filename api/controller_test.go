package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	//recapi "github.com/James-Trauger/Recipouir/api"
	"github.com/James-Trauger/Recipouir/model"
)

func validLoginResponse(expected model.Login, actual model.User, expErr, actualErr error, expStatus, actualStatus int) bool {
	return expErr == actualErr && expected.EqualUser(actual) && expStatus == actualStatus
}

func TestSignup(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// create a new request with a raw string of a login
	req := httptest.NewRequest(http.MethodPost, "/api/signup",
		strings.NewReader(`{"uname":"james","pass":"hello"}`))
	loginResponse, status, err := Signup(req, ctx)

	if loginResponse == nil || !validLoginResponse(model.Login{Uname: "james", Pass: "hello"}, *loginResponse,
		nil, err, http.StatusOK, status) {

		t.Fatal(err)
	} else {
		// delete the added user
		err = DeleteUser(loginResponse, ctx)
		if err != nil {
			t.Log("couldn't delete added user -> " + err.Error())
		}
	}

	// test a request using the Login struct
	login := model.Login{
		Uname: "trau",
		Pass:  "world",
	}
	loginJSON, err := json.Marshal(&login)
	if err != nil {
		t.Fatal(err)
	}
	// create the request and pass it to the controller
	req = httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewReader(loginJSON))

	loginResponse, status, err = Signup(req, ctx)
	if loginResponse == nil || !validLoginResponse(login, *loginResponse, nil, err, http.StatusOK, status) {
		t.Fatal(err)
	} else {
		// delete the added user
		err = DeleteUser(loginResponse, ctx)
		if err != nil {
			t.Log("couldn't delete added user -> " + err.Error())
		}
	}
}
