package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

var ings = []Ingredient{
	NewIng("flour", 1, 2, "cup"),
	NewIng("vanilla", 5, 1, "gram"),
	NewIng("chocolate", 2, 1, "oz"),
}

var ingsJson = []string{
	`{"Name":"flour","Amount":"1/2","Unit":"cup"}`,
	`{"Name":"vanilla","Amount":"5","Unit":"gram"}`,
	`{"Name":"chocolate","Amount":"2","Unit":"oz"}`,
}

var steps = []string{
	"mix the flour",
	"add the vanilla",
}

func TestIngMarshal(t *testing.T) {
	/*ings := []Ingredient{
		NewIngredient("flour", "cup", big.NewRat(1, 2)),
	}*/

	for i, js := range ings {
		act, err := json.Marshal(js)
		if err != nil {
			fmt.Println("error marashaling")
			t.Fail()
		}
		if string(act) != ingsJson[i] {
			t.Fatalf("ingredient was not marshalled correctly\nReceived: %s\nexpected: %s", string(act), ingsJson[i])
		}
	}
}

func TestIngUnmarshal(t *testing.T) {
	for _, ingredient := range ings {
		// marshall the ingredient
		jsonBytes, err := json.Marshal(ingredient)
		if err != nil {
			fmt.Println("error marashaling")
			t.Fail()
		}

		var marshIng Ingredient

		// unmarshall
		err = json.Unmarshal(jsonBytes, &marshIng)
		if err != nil {
			t.Fatal(err)
		}

		// the fields need to exactly match
		if ingredient.Name != marshIng.Name {
			t.Fatalf("name of unmarshalled ingredient is incorrect\nExpected: %s\nReceived: %s",
				ingredient.Name, marshIng.Name)
		}
		if ingredient.Amount.Cmp(&marshIng.Amount.Rat) != 0 {
			t.Fatalf("amount of unmarshalled ingredient is incorrect\nExpected: %s\nReceived: %s",
				ingredient.Amount, marshIng.Amount)
		}
		if ingredient.Unit != marshIng.Unit {
			t.Fatalf("amount of unmarshalled ingredient is incorrect\nExpected: %s\nReceived: %s",
				ingredient.Unit, marshIng.Unit)
		}
	}
}

func TestRecipeMarshal(t *testing.T) {
	recipe := NewRecipe("mine", "ned", ings, steps)
	raw, _ := json.RawMessage(fmt.Sprintf("{\"Name\":\"mine\",\"Ingredients\":[%s],\"Steps\":[\"mix the flour\",\"add the vanilla\"],\"CreatedBy\":\"ned\"}", strings.Join(ingsJson, ","))).MarshalJSON()
	//raw, _ := json.RawMessage(`{"name":"mine","ingredients":[{"name":"flour","amount":"1/2","unit":"cup"},{"name":"vanilla","amount":"5","unit":"gram"}],"steps":["mix the flour","add the vanilla"],"username":"ned"}`).MarshalJSON()

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		fmt.Println("error marashaling")
		t.Fail()
	}

	if !bytes.Equal(jsonBytes, raw) {
		t.Fatalf("marshalling a Recipe faild\nExpected: %s\nReceived: %s", string(raw), string(jsonBytes))
	}
}

func TestRecipeUnmarshal(t *testing.T) {
	recipe := NewRecipe("mine", "ned", ings, steps)

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		fmt.Println("error marshalling Recipe")
		t.Fail()
	}

	var unmarshRecipe Recipe
	if json.Unmarshal(jsonBytes, &unmarshRecipe) != nil {
		t.Fatal("error unmarshalling recipe")
	}

	// compare names
	if recipe.Name != unmarshRecipe.Name {
		t.Fatalf("name of recipe unmarshalled incorrectly\nExpected: %s\nReceived: %s\n", recipe.Name, unmarshRecipe.Name)
	}

	// compare ingredients
	if !EqualIngredients(recipe.Ingredients, unmarshRecipe.Ingredients) {
		t.Fatalf("ingredients of recipes unmarshalled incorrectly\nExpected: %s\nReceived: %s\n", recipe.Ingredients, unmarshRecipe.Ingredients)
	}

	//compare steps
	if !slices.Equal(recipe.Steps, unmarshRecipe.Steps) {
		t.Fatalf("ingredients of recipes unmarshalled incorrectly\nExpected: %s\n Received: %s\n", recipe.Steps, unmarshRecipe.Steps)
	}

	if recipe.CreatedBy != unmarshRecipe.CreatedBy {
		t.Fatalf("name of recipe creator unmarshalled incorrectly\nExpected: %s\nReceived: %s\n", recipe.CreatedBy, unmarshRecipe.CreatedBy)
	}
}

func TestIngredientBSON(t *testing.T) {
	for _, ing := range ings {
		// encode to bson
		bs, err := bson.Marshal(&ing)
		if err != nil {
			t.Fatal(err)
		}

		// decode back to original type
		var unmarshalledIng Ingredient
		if err = bson.Unmarshal(bs, &unmarshalledIng); err != nil {
			t.Fatal(err)
		}
		// compare to original value
		if !ing.Equal(&unmarshalledIng) {
			t.Fatalf("ingredient unmarshalled bson incorrectly\nExpected:%v\nReceived:%v", ing, unmarshalledIng)
		}
	}
}

func TestRecipeBSON(t *testing.T) {
	recipe := NewRecipe("mine", "ned", ings, steps)

	bsonBytes, err := bson.Marshal(&recipe)
	if err != nil {
		t.Fatal(err)
	}

	var unmarshalledRecipe Recipe
	if err = bson.Unmarshal(bsonBytes, &unmarshalledRecipe); err != nil {
		t.Fatal(err)
	}

	if !recipe.Equal(&unmarshalledRecipe) {
		t.Fatalf("unmarshalled recipe does not match original recipe\nExpected: %v\nReceived: %v\n", recipe, unmarshalledRecipe)
	}
}
