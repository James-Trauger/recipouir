package utils

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
type RSA2048SigningMethod struct {
	Hash crypto.Hash
}

func (sm *RSA2048SigningMethod) Alg() string { return "RS2048" }

func (sm *RSA2048SigningMethod) Sign(signingString string, key any) ([]byte, error) {
	var rsaKey *rsa.PrivateKey
	var ok bool
	// needs to be an rsa key
	if rsaKey, ok = key.(*rsa.PrivateKey); !ok {
		return nil, jwt.ErrInvalidKeyType
	}

	if !sm.Hash.Available() {
		return nil, jwt.ErrHashUnavailable
	}
	// hash the raw token
	hashed := sm.Hash.New()
	hashed.Write([]byte(signingString))
	// sign the hashed token
	if sigBytes, err := rsa.SignPKCS1v15(nil, rsaKey, sm.Hash, nil); err == nil {
		return sigBytes, nil
	} else {
		return nil, err
	}
}

func (sm *RSA2048SigningMethod) Verify(signingString string, sig []byte, key any) error {
	var rsaKey *rsa.PublicKey
	var ok bool
	// needs to be an rsa key
	if rsaKey, ok = key.(*rsa.PublicKey); !ok {
		return jwt.ErrInvalidKeyType
	}

	if !sm.Hash.Available() {
		return jwt.ErrHashUnavailable
	}
	// hash the string
	hasher := sm.Hash.New()
	hasher.Write([]byte(signingString))
	return rsa.VerifyPKCS1v15(rsaKey, sm.Hash, hasher.Sum(nil), sig)
}*/

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (uc UserClaims) Validate() error {
	if uc.Username == "" {
		return errors.New("empty username")
	}
	expires, err := uc.GetExpirationTime()
	if err != nil {
		return err
	}
	if expires.Compare(time.Now()) < 0 {
		return jwt.ErrTokenExpired
	}
	return nil
}

func NewToken(user string) (string, error) {
	// return a jwt token using RSA, expires a day from now
	t := jwt.New(jwt.SigningMethodRS256)
	// set the claims
	t.Claims = &UserClaims{
		user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	return t.SignedString(PrivateKey)
}

func VerifyToken(t *jwt.Token) (any, error) { return PublicKey, nil }

// returns the token from an http header (does NOT validate the token)
func ParseTokenFromHeader(head *http.Header) (string, error) {
	// token is in the Authorizaiton header
	authHeader := strings.Split(head.Get("Authorization"), " ")
	// header is in the form `Bearer [token]`
	if len(authHeader) != 2 || authHeader[0] != "Bearer" {
		return "", errors.New("malformed authorization header, expected \"Authorization: Bearer [token]\"") // invalid header
	}

	return authHeader[1], nil
}

/*
accepts a request header to parse and validate a jwt token
*/
func ValidToken(rawToken string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(rawToken, &UserClaims{}, VerifyToken)
	if err != nil {
		//JSONError(w, http.StatusInternalServerError, errors.New("couldn't parse token -> "+err.Error()))
		return nil, err
	}
	claims := token.Claims.(*UserClaims)
	return claims, claims.Validate()
}
