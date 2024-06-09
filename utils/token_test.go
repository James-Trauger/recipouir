package utils_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/James-Trauger/Recipouir/utils"
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
	tokenString, err := utils.NewToken(uname)
	if err != nil {
		t.Fatal(err)
	}
	// populate the header
	req.Header.Add("Authorization", "Bearer "+tokenString)

	// parse the header
	rawToken, err := utils.ParseTokenFromHeader(&req.Header)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := utils.ValidToken(rawToken)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("claims: %v\n", claims)
	dur := time.Until(claims.ExpiresAt.Time)
	if claims.Username != "admin" && dur < (24*time.Hour) && dur > ((23*time.Hour)+(59*time.Minute)) {
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
