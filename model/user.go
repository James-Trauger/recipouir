package model

import (
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
	}
}

/*
used when making login request for a user
Either Uname or Email must be present, or both
*/
type Login struct {
	Uname *string `json:"uname"`
	Pass  *string `json:"pass"`
}

func HashPassword(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), HashCost)
}
