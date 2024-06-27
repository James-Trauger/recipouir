package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	//recapi "github.com/James-Trauger/Recipouir/api"
	"github.com/James-Trauger/Recipouir/model"
)

func validLoginResponse(expected model.Login, actual model.User, expErr, actualErr error, expStatus, actualStatus int) bool {
	return expErr == actualErr && expected.EqualUser(actual) && expStatus == actualStatus
}

func TestSignup(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// create a new request with a raw string of a login
	req := httptest.NewRequest(http.MethodPost, "/api/signup",
		strings.NewReader(`{"uname":"james","pass":"hello"}`))
	loginResponse, status, err := Signup(req, ctx)

	if loginResponse == nil || !validLoginResponse(model.Login{Uname: "james", Pass: "hello"}, *loginResponse,
		nil, err, http.StatusOK, status) {

		t.Fatal(err)
	} else {
		// delete the added user
		err = DeleteUser(loginResponse, ctx)
		if err != nil {
			t.Log("couldn't delete added user -> " + err.Error())
		}
	}

	// test a request using the Login struct
	login := model.Login{
		Uname: "trau",
		Pass:  "world",
	}
	loginJSON, err := json.Marshal(&login)
	if err != nil {
		t.Fatal(err)
	}
	// create the request and pass it to the controller
	req = httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewReader(loginJSON))

	loginResponse, status, err = Signup(req, ctx)
	if loginResponse == nil || !validLoginResponse(login, *loginResponse, nil, err, http.StatusOK, status) {
		t.Fatal(err)
	} else {
		// delete the added user
		err = DeleteUser(loginResponse, ctx)
		if err != nil {
			t.Log("couldn't delete added user -> " + err.Error())
		}
	}
}

func TestInsertGetRecipe(t *testing.T) {
	// create a new recipe
	user := "ned"
	rec := model.NewRecipe("cookies", user, []model.Ingredient{model.NewIng("flour", 2, 1, "cup")},
		[]string{"mix flour, sugar, and milk"})
	// insert the recipe
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := InsertRecipe(*rec, user, ctx)
	if err != nil {
		t.Fatal(err)
	}

	//delete the inserted recipe
	defer func(recipe, user string, ctx context.Context) {
		// delete from db
		err = DeleteRecipe(rec.Name, user, ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(rec.Name, user, ctx)

	// retrieve the recipe
	retreivedRecipe, err := GetRecipe(user, rec.Name, ctx)
	if err != nil {
		t.Fatal(err)
	}

	// copmare recipes
	if !rec.Equal(retreivedRecipe) {
		t.Fatalf("retrieved recipe is not the same as the inserted one\nInserted: %v\nRetreived: %v", rec, retreivedRecipe)
	}
}

func TestInsertManyRecipe(t *testing.T) {
	user := "ned"
	recs := []model.Recipe{
		*model.NewRecipe("cookies", user, []model.Ingredient{model.NewIng("flour", 2, 1, "cup")},
			[]string{"mix flour, sugar, and milk"}),
		*model.NewRecipe("brownies", user, []model.Ingredient{model.NewIng("sugar", 1, 3, "cup"), model.NewIng("butter", 1, 1, "stick")},
			[]string{"combine sugar and butter"}),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// insert the recipes
	for _, rec := range recs {
		if err := InsertRecipe(rec, user, ctx); err != nil {
			t.Fatal(err)
		}
	}

	// retrieve the recipes and delete them
	for _, rec := range recs {
		retRec, err := GetRecipe(user, rec.Name, ctx)
		if err != nil {
			t.Error(err)
		}

		/* delete the recipe
		if err = DeleteRecipe(rec.Name, user, ctx); err != nil {
			t.Error(err)
		}*/

		// compare the original recipe with the retrieved recipes
		if !rec.Equal(retRec) {
			t.Errorf("Inserted recipe is not the same as the original\nExpected: %v\n,Received: %v\n", rec, retRec)
		}
	}

}
