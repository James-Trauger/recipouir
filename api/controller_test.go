package api

import (
	"context"
	"testing"
	"time"

	//recapi "github.com/James-Trauger/Recipouir/api"
	"github.com/James-Trauger/Recipouir/model"
)

func validLoginResponse(expected model.Login, actual model.User, expErr, actualErr error) bool {
	return expErr == actualErr && expected.EqualUser(actual)
}

func TestSignup(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	login := model.Login{
		Uname: "james",
		Pass:  "hello",
	}

	loginResponse, err := Signup(&login, ctx)

	if loginResponse == nil || !validLoginResponse(model.Login{Uname: "james", Pass: "hello"}, *loginResponse,
		nil, err) {

		t.Fatal(err)
	} else {
		// delete the added user
		err = DeleteUser(loginResponse, ctx)
		if err != nil {
			t.Log("couldn't delete added user -> " + err.Error())
		}
	}
}

// test login

func TestInsertGetRecipe(t *testing.T) {
	// create a new recipe
	user := "ned"
	rec := model.NewRecipe("oatmeal cookies", user, []model.Ingredient{model.NewIng("oatmeal", 1, 2, "cup")},
		[]string{"combine oats and sugar"})
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
		*model.NewRecipe("oatmeal cookies", user, []model.Ingredient{model.NewIng("oatmeal", 1, 2, "cup")},
			[]string{"combine oats and sugar"}),
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

		// delete the recipe
		if err = DeleteRecipe(rec.Name, user, ctx); err != nil {
			t.Error(err)
		}

		// compare the original recipe with the retrieved recipes
		if !rec.Equal(retRec) {
			t.Errorf("Inserted recipe is not the same as the original\nExpected: %v\n,Received: %v\n", rec, retRec)
		}
	}

}
