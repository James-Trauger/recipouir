package reciauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidToken(t *testing.T) {
	// test username
	uname := "admin"
	// empty request
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// generate the token
	tokenString, err := NewToken(uname)
	if err != nil {
		t.Fatal(err)
	}
	// populate the header
	//req.Header.Add("Authorization", "Bearer "+tokenString)
	AddTokenHeader(req, tokenString)

	// parse the header
	rawToken, err := ParseTokenFromHeader(&req.Header)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ValidToken(rawToken)
	if err != nil {
		t.Fatal(err)
	}
	dur := time.Until(claims.ExpiresAt.Time)
	if claims.Username != "admin" || dur < 0 || claims.Issuer != "Recipouir" {
		t.Fatalf("invalid token, incorrect username or expiry date\nusername=%s\nexpires at %s\n", claims.Username, claims.ExpiresAt)
	}
}

func TestInvalidToken(t *testing.T) {
	falseToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiaXNzIjoiUmVjaXBvdWlyIiwiZXhwIjoxNzIyMzg2NTc0fQ.DuD8wsN6GZTjAqhVIS_N2hTTZZn3pum-L4cxHYQ-bihXllK-1A9DQR22UbN87_en5aRO-bH0sE8RaXzkJvVC_10AfsEYL0nd7CCYF-Ir-Fe0h-xN-xBAtI3fcs6vgkL2Dmc2uL3EXA9-VY3zguP8cPH74FBODft5kqmxeSClJuBZvCEGlX1mSIguaIaOponYNCneap7pNmg8lE6L777TjO2i78BppnfCbMP2r22QLZy31auzOeLqtKuVBqQ6GfKF8pKh3zNd7LEZ2z5ybMaHg7kztRpFFHxAIL58t-NWYYLXS75xw19-tbOEUCNKSyNcX9uMnjP-dLaUWoZFuSMgcg"
	_, err := ValidToken(falseToken)
	if err == nil {
		t.Fatal("invalid token was successfully verified")
	}
	if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		fmt.Println(err)
	}
}

func TestPath(t *testing.T) {
	u, err := url.Parse("/api/user/james")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u.Path)
}
