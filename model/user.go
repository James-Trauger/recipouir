package model

import (
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// cost of hashing passwords with bcrypt
const HashCost = 12

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	First    string             `json:"first" validate:"required, min=2, max=50"`
	Last     string             `json:"last" validate:"required, min=2, max=50"`
	Pass     *[]byte            `json:"password" validate:"required, min=8, max=64"`
	Email    string             `json:"email" validate:"email, required"`
	Token    *string            `json:"token"`
	UserType string             `json:"UserType" validate:"required, eq=ADMIN|eq=USER"`
	Refresh  *string            `json:"refresh"`
	Uname    string             `json:"username"`
	UserID   *uuid.UUID         `json:"userid"` // random, unguessable, and unique value
}

func NewUser(first, last, uname, pass, email string) *User {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(pass), HashCost)
	if err != nil {
		log.Println("couldn't generate password, " + err.Error())
		return nil
	}
	// TODO make sure the user doesn't already exist
	// generate token and refresh

	return &User{
		First:    first,
		Last:     last,
		Pass:     &passBytes,
		Email:    email,
		Uname:    uname,
		UserType: "USER",
	}
}

/*
used when making login request for a user
Either Uname or Email must be present, or both
*/
type Login struct {
	Uname *string `json:"uname"`
	Email *string `json:"email"`
	Pass  *string `json:"pass"`
}
