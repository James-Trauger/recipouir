package reciauth

import (
	"context"

	db "github.com/James-Trauger/Recipouir/database"
	"github.com/James-Trauger/Recipouir/model"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

/* pass a user's login credentials and retrieve a User if it matches their info in the db */
func Authenticate(login *model.Login, ctx context.Context) (*model.User, error) {
	filter := bson.M{"_id": login.Uname}
	result := db.OpenCollection(db.Client, db.DbName, "user").FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return nil, err
	}

	var user model.User
	err := result.Decode(&user)
	// compare passwords
	if isAuthenticated := bcrypt.CompareHashAndPassword(user.Pass.Data, []byte(login.Pass)); isAuthenticated != nil {
		return nil, isAuthenticated
	}

	if err != nil {
		// couldn't decod the user
		//return errors.New("internal server error")
		return nil, err
	} else {
		return &user, nil
	}
}

func ValidPassword(pass *[]byte) bool {
	return false
}
