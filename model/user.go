package model

import (
	"encoding/json"
	"io"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// cost of hashing passwords with bcrypt
const HashCost = 12

type User struct {
	//ID       primitive.ObjectID `bson:"_id"`
	// primary key in the db, the same as the username
	ID string `bson:"_id"`
	// hashed password
	Pass     *primitive.Binary `bson:"password" validate:"required, min=8, max=64"`
	UserType string            `json:"userType" validate:"required, eq=ADMIN|eq=USER"`
	Username string            `json:"username"`
}

func NewUser(uname, pass string) *User {
	hashed, err := HashPassword(pass)
	if err != nil {
		log.Println("couldn't generate password, " + err.Error())
		return nil
	}
	// TODO make sure the user doesn't already exist
	// generate token and refresh

	return &User{
		Username: uname,
		Pass: &primitive.Binary{
			Data: hashed,
		},
		UserType: "USER",
		//ID:       primitive.NewObjectID(),
		ID: uname,
	}
}

func (u *User) Equal(u2 User) bool {
	return u.Username == u2.Username && u.Pass.Equal(*u2.Pass) &&
		u.UserType == u2.UserType && u.ID == u2.ID
}

/*
used when making a login request for a user
*/
type Login struct {
	Uname string `json:"uname"`
	Pass  string `json:"pass"`
}

func (l *Login) Equal(l2 Login) bool {
	return l != nil && l.Uname == l2.Uname && l.Pass == l2.Pass
}

func (l *Login) EqualUser(u User) bool {
	return l != nil && l.Uname == u.Username && bcrypt.CompareHashAndPassword(u.Pass.Data, []byte(l.Pass)) == nil
}

func HashPassword(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), HashCost)
}

// returns a user in the body of a request
func ExtractLogin(r io.ReadCloser) (*Login, error) {
	// wrap the request body so it can be decoded
	var l Login
	buf := json.NewDecoder(r)
	err := buf.Decode(&l)
	return &l, err
}
