package reciauth

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
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
	req.Header.Add("Authorization", "Bearer "+tokenString)

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

func TestPath(t *testing.T) {
	u, err := url.Parse("/api/user/james")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u.Path)
}
