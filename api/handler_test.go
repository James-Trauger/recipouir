package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	reciauth "github.com/James-Trauger/Recipouir/auth"
	"github.com/James-Trauger/Recipouir/model"
	"github.com/James-Trauger/Recipouir/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	signupPath = "/api/signup"
	loginPath  = "/api/login"
)

// example users
var (
	ned = model.NewLogin("ned", "honor")
	jon = model.NewLogin("jon", "snow")
)

func TestDefaultHandle(t *testing.T) {
	serv := httptest.NewServer(RootHandler())

	req := httptest.NewRequest(http.MethodGet, serv.URL, nil)
	w := httptest.NewRecorder()
	RootHandler().ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode) // GET
	}
	utils.DrainClose(resp.Body)

	// test method not found
	req = httptest.NewRequest(http.MethodPut, serv.URL, nil)
	w = httptest.NewRecorder()
	RootHandler().ServeHTTP(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		log.Fatal("incorrect response code for an unsupported method")
	}
	utils.DrainClose(resp.Body)

	// test Options method
	req = httptest.NewRequest(http.MethodOptions, serv.URL, nil)
	w = httptest.NewRecorder()
	RootHandler().ServeHTTP(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("options method not supported")
	}
	//fmt.Println(resp.Header.Get("Allow"))
	utils.DrainClose(resp.Body)
}

func loginReader(l *model.Login) (*bytes.Reader, error) {
	jsonBytes, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonBytes), nil
}

func TestLoginHandler(t *testing.T) {
	reader, err := loginReader(ned)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := httptest.NewRequest(http.MethodPost, loginPath, reader).WithContext(ctx)
	w := httptest.NewRecorder()
	HandleLogin().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatal("incorrect status code ", w.Code)
	}

	// verify the jwt token
	var tokResp tokenResponse
	err = json.NewDecoder(w.Body).Decode(&tokResp)
	if err != nil {
		t.Fatal("couldn't decode json response ", err)
	}

	claims, err := reciauth.ValidToken(tokResp.Token)
	if err != nil {
		t.Fatal(err)
	}

	if err = claims.Validate(); err != nil {
		t.Fatal("invalid token ", err)
	}
}

func TestInvalidLoginHandler(t *testing.T) {
	creds := model.NewLogin(ned.Uname, "ambition")
	reader, err := loginReader(creds)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := httptest.NewRequest(http.MethodPost, loginPath, reader).WithContext(ctx)
	w := httptest.NewRecorder()
	HandleLogin().ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatal("incorrect status code ", w.Code)
	}

	//////////////////////////
	io.Copy(os.Stdout, w.Body)

}

func TestSignupHandler(t *testing.T) {

	// create the signup information and encode it to json
	reader, err := loginReader(jon)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	signupReq := httptest.NewRequest(http.MethodPost, signupPath, reader).WithContext(ctx)
	w := httptest.NewRecorder()

	SignupHandler().ServeHTTP(w, signupReq)

	// check the user was created
	user, err := Login(jon, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusOK {
		t.Fatalf("request was not accepted: %d\n", w.Code)
	}

	// verify username
	if user.Username != jon.Uname {
		t.Fatalf("usernames do not match\nExpected:%s\nReceived:%s", jon.Uname, user.Username)
	}
	// verify password
	if bcrypt.CompareHashAndPassword(user.Pass.Data, []byte(jon.Pass)) != nil {
		t.Fatal("passwords do not match")
	}

	//delete the user
	if err = DeleteUser(user, ctx); err != nil {
		t.Fatal(err)
	}
}

// signup a user that already exists
func TestSignupDuplicateHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ensure the user already exists
	if _, err := Login(ned, ctx); err != nil {
		t.Fatal("user may not exist: ", err)
	}

	reader, err := loginReader(ned)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, signupPath, reader).WithContext(ctx)

	SignupHandler().ServeHTTP(w, r)

	// check response code
	if w.Code != http.StatusConflict {
		t.Fatal("duplicate user was created")
	}
}

func TestTokenJson(t *testing.T) {
	tk := tokenResponse{
		Token: "signedToken",
	}
	bts, err := json.Marshal(&tk)
	if err != nil {
		t.Fatal(err)
	}

	expectedBts := []byte("{\"token\":\"signedToken\"}")
	if !bytes.Equal(bts, expectedBts) {
		t.Fatal("json does not match")
	}
}
