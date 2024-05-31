package utils

import (
	"math/big"
)

type Ingredient struct {
	Name   string   `json:"name"`
	Amount *big.Rat `json:"amount"` // fraction
	Unit   string   `json:"unit"`   // cup? grams? ml?
}

/* returns an ingredient with a passeed Rational number as the amount */
func NewIngredient(n string, amt *big.Rat, un string) Ingredient {
	return Ingredient{
		Name:   n,
		Amount: amt,
		Unit:   un,
	}
}

/* returns an ingredient with a numerator and denominator as the amount */
func NewIng(n string, num int64, denom int64, un string) Ingredient {
	return Ingredient{
		Name:   n,
		Amount: big.NewRat(num, denom),
		Unit:   un,
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

type Recipe struct {
	Name  string       `json:"name"`
	Ings  []Ingredient `json:"ingredients"`
	Steps []string     `json:"steps"`
}

/*
creates a new recipe based on the name, ingredients, and steps.
pass nil to the ingredients and sps argument if you want to add them yourself
*/
func NewRecipe(n string, ingredients []Ingredient, sps []string) *Recipe {
	return &Recipe{
		Name:  n,
		Ings:  ingredients,
		Steps: sps,
	}
}

func (r *Recipe) AddIngredient(i ...Ingredient) {
	r.Ings = append(r.Ings, i...)
}

func (r *Recipe) AddStep(s ...string) {
	r.Steps = append(r.Steps, s...)
}
