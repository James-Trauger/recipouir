package model

import (
	"errors"
	"fmt"
	"math/big"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
)

type Ingredient struct {
	Name   string   `json:"name"`
	Amount *big.Rat `json:"amount"` // fraction
	Unit   string   `json:"unit"`   // cup? grams? ml?
}

/* returns an ingredient with a passeed Rational number as the amount */
func NewIngredient(name, unit string, amt *big.Rat) Ingredient {
	return Ingredient{
		Name:   name,
		Amount: amt,
		Unit:   unit,
	}
}

/* returns an ingredient with a numerator and denominator as the amount */
func NewIng(name string, num int64, denom int64, unit string) Ingredient {
	return Ingredient{
		Name:   name,
		Amount: big.NewRat(num, denom),
		Unit:   unit,
	}
}

/* equal ingredients must have the same exact fields */
func (i1 *Ingredient) Equal(i2 *Ingredient) bool {
	return i1.Name == i2.Name && i1.Amount.Cmp(i2.Amount) == 0 && i1.Unit == i2.Unit
}

/* a list of ingredients must have the same size and same fields */
func EqualIngredients(i1, i2 []Ingredient) bool {
	if len(i1) != len(i2) {
		return false
	}
	for i, ing1 := range i1 {
		if !ing1.Equal(&i2[i]) {
			return false
		}
	}
	return true
}

func (i *Ingredient) MarshalBSON() ([]byte, error) {
	document := bson.D{
		{"name", i.Name},
		{"amount", i.Amount},
		{"unit", i.Unit},
	}
	if i.Amount == nil || i.Name == "" || i.Unit == "" {
		return nil, errors.New("invalid ingredient")
	}
	if i.Amount.Denom().Cmp(big.NewInt(1)) == 0 {
		return []byte(fmt.Sprintf("{ name: '%s', amount: %s, unit: '%s' }", i.Name, i.Amount.Num(), i.Unit)), nil
	} else {
		return []byte(fmt.Sprintf("{ name: '%s', amount: %s, unit: '%s' }", i.Name, i.Amount, i.Unit)), nil
	}
}

type Recipe struct {
	Name      string       `json:"name"`
	Ings      []Ingredient `json:"ingredients"`
	Steps     []string     `json:"steps"`
	CreatedBy string       `json:"username"`
}

/*
creates a new recipe based on the name, ingredients, and steps.
pass nil to the ingredients and sps argument if you want to add them yourself
*/
func NewRecipe(name, creator string, ingredients []Ingredient, steps []string) *Recipe {
	return &Recipe{
		Name:      name,
		Ings:      ingredients,
		Steps:     steps,
		CreatedBy: creator,
	}
}

func (r *Recipe) Equal(r2 *Recipe) bool {
	return r.Name == r2.Name && EqualIngredients(r.Ings, r2.Ings) &&
		slices.Equal(r.Steps, r2.Steps) && r.CreatedBy == r2.CreatedBy
}

func (r *Recipe) AddIngredient(i ...Ingredient) {
	r.Ings = append(r.Ings, i...)
}

func (r *Recipe) AddStep(s ...string) {
	r.Steps = append(r.Steps, s...)
}
