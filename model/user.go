package model

import (
	"bytes"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// cost of hashing passwords with bcrypt
const HashCost = 12

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Pass     *[]byte            `json:"password" validate:"required, min=8, max=64"`
	UserType string             `json:"userType" validate:"required, eq=ADMIN|eq=USER"`
	Username string             `json:"username"`
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
		Pass:     &hashed,
		UserType: "USER",
		ID:       primitive.NewObjectID(),
	}
}

func (u *User) Equal(u2 User) bool {
	return u.Username == u2.Username && bytes.Equal(*u.Pass, *u2.Pass) &&
		u.UserType == u2.UserType && u.ID.Hex() == u2.ID.Hex()
}

/*
used when making login request for a user
Either Uname or Email must be present, or both
*/
type Login struct {
	Uname string `json:"uname"`
	Pass  string `json:"pass"`
}

func (l *Login) Equal(l2 Login) bool {
	return l != nil && l.Uname == l2.Uname && l.Pass == l2.Pass
}

func (l *Login) EqualUser(u User) bool {
	return l != nil && l.Uname == u.Username && bcrypt.CompareHashAndPassword(*u.Pass, []byte(l.Pass)) == nil
}

func HashPassword(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), HashCost)
}
