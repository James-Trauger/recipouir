package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"testing"
)

var ings = []Ingredient{
	NewIng("flour", 1, 2, "cup"),
	NewIng("vanilla", 5, 1, "gram"),
}

var steps = []string{
	"mix the flour",
	"add the vanilla",
}

func TestIngMarshal(t *testing.T) {
	// strings of the ingredients
	ingsJson := []string{
		"{\"name\":\"flour\",\"amount\":\"1/2\",\"unit\":\"cup\"}",
	}
	ings := []Ingredient{
		NewIngredient("flour", big.NewRat(1, 2), "cup"),
	}

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
		json.Unmarshal(jsonBytes, &marshIng)

		// the fields need to exactly match
		if ingredient.Name != marshIng.Name {
			t.Fatalf("name of unmarshalled ingredient is incorrect\nExpected: %s\nReceived: %s",
				ingredient.Name, marshIng.Name)
		}
		if ingredient.Amount.Cmp(marshIng.Amount) != 0 {
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
	recipe := NewRecipe("mine", ings, steps)
	raw, _ := json.RawMessage(`{"name":"mine","ingredients":[{"name":"flour","amount":"1/2","unit":"cup"},{"name":"vanilla","amount":"5","unit":"gram"}],"steps":["mix the flour","add the vanilla"]}`).MarshalJSON()

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
	recipe := NewRecipe("mine", ings, steps)

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
	if !EqualIngredients(recipe.Ings, unmarshRecipe.Ings) {
		t.Fatalf("ingredients of recipes unmarshalled incorrectly\nExpected: %s\nReceived: %s\n", recipe.Ings, unmarshRecipe.Ings)
	}

	//compare steps
	if !slices.Equal(recipe.Steps, unmarshRecipe.Steps) {
		t.Fatalf("ingredients of recipes unmarshalled incorrectly\nExpected: %s\n Received: %s\n", recipe.Steps, unmarshRecipe.Steps)
	}
}
