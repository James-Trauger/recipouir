package reciauth

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

var (
	ErrTokenMissing = errors.New("no token was found in the authorization header, attach the token or login to receive a new one")
)

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// returns nil on success
// validates the token's claim by verifying a username exists and the token is not expired
func (uc UserClaims) Validate() error {
	if uc.Username == "" {
		return errors.New("empty username")
	}
	expires, err := uc.GetExpirationTime()
	if err != nil {
		return err
	}
	// past the expiry date
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
			Issuer:    "Recipouir",
		},
	}
	return t.SignedString(PrivateKey)
}

func VerifyToken(t *jwt.Token) (any, error) { return PublicKey, nil }

// returns the token from an http header (does NOT validate the token)
func ParseTokenFromHeader(head *http.Header) (string, error) {
	authHeader := head.Get("Authorization")
	if authHeader == "" {
		return "", ErrTokenMissing
	}

	// token is in the Authorizaiton header
	bearer := strings.Split(head.Get("Authorization"), " ")

	// header is in the form `Bearer [token]`
	if len(bearer) != 2 || bearer[0] != "Bearer" {
		return "", jwt.ErrTokenMalformed // invalid header
	}

	return bearer[1], nil
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

/*
authorizes a user given a raw token and a username
*/
func Authroize(token, user string) bool {
	// validate the token and retrieve the claims
	claims, err := ValidToken(token)
	if err != nil {
		return false
	}
	// username of the token must match the passed username
	return claims.Username == user
}

func AddTokenHeader(r *http.Request, token string) {
	r.Header.Add("Authorization", "Bearer "+token)
}
