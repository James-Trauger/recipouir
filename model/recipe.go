package model

import (
	"errors"
	"math/big"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type Fraction struct {
	big.Rat `json:"Amount"`
}

/*
func (f *Fraction) MarshalBSONValue() (bsontype.Type, []byte, error) {
	var err error

	// initialize buffer with type string
	buf := bytes.NewBuffer([]byte{0x02})
	// add field name
	if _, err = buf.Write([]byte("amount")); err != nil {
		return bson.TypeUndefined, nil, err
	}
	// add terminal
	if err = buf.WriteByte(0x0); err != nil {
		return bson.TypeUndefined, nil, err
	}

	// retrieve binary representation of the field value
	fieldValue, err := f.MarshalText()
	if err != nil {
		return bson.TypeUndefined, nil, err
	}

	// add the length + 1 including null character
	fieldLen := len(fieldValue)
	if err = binary.Write(buf, binary.BigEndian, int32(fieldLen+1)); err != nil {
		return bson.TypeUndefined, nil, err
	}

	// add the field value
	if _, err = buf.Write(fieldValue); err != nil {
		return bson.TypeUndefined, nil, errors.New("couldn't encode field name, amount for Fraction")
	}
	// add terminal
	if err = buf.WriteByte(0x0); err != nil {
		return bson.TypeUndefined, nil, err
	}

	return bson.TypeString, nil, err
}*/

func (f *Fraction) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if f.Denom().Cmp(big.NewInt(1)) == 0 {
		// whole number
		return bson.MarshalValue(f.Num().String())
	} else {
		return bson.MarshalValue(f.String())
	}
}

// use the UnmarshalText function of big.Rat to decode it from a string
func (f *Fraction) UnmarshalBSONValue(t bsontype.Type, data []byte) error {

	if t != bson.TypeString {
		return errors.New("Fraction must be encoded as a bson string")
	}

	var s string
	if err := bson.UnmarshalValue(t, data, &s); err != nil {
		return err
	}

	if err := f.UnmarshalText([]byte(s)); err != nil {
		return err
	}

	return nil
}

type Ingredient struct {
	Name   string    `json:"Name"`
	Amount *Fraction `bson:"amount"` // fraction
	Unit   string    `json:"Unit"`   // cup? grams? ml?
}

/* returns an ingredient with a passeed Rational number as the amount */
func NewIngredient(name, unit string, amt *Fraction) Ingredient {
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
		Amount: &Fraction{*big.NewRat(num, denom)},
		Unit:   unit,
	}
}

/* equal ingredients must have the same exact fields */
func (i1 *Ingredient) Equal(i2 *Ingredient) bool {
	return i1.Name == i2.Name && i1.Amount.Cmp(&i2.Amount.Rat) == 0 && i1.Unit == i2.Unit
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
	Name        string       `json:"Name"`
	Ingredients []Ingredient `json:"Ingredients"`
	Steps       []string     `json:"Steps"`
	CreatedBy   string       `json:"CreatedBy"`
}

/*
creates a new recipe based on the name, ingredients, and steps.
pass nil to the ingredients and sps argument if you want to add them yourself
*/
func NewRecipe(name, creator string, ingredients []Ingredient, steps []string) *Recipe {
	return &Recipe{
		Name:        name,
		Ingredients: ingredients,
		Steps:       steps,
		CreatedBy:   creator,
	}
}

func (r *Recipe) Equal(r2 *Recipe) bool {
	return r.Name == r2.Name && EqualIngredients(r.Ingredients, r2.Ingredients) &&
		slices.Equal(r.Steps, r2.Steps) && r.CreatedBy == r2.CreatedBy
}

func (r *Recipe) AddIngredient(i ...Ingredient) {
	r.Ingredients = append(r.Ingredients, i...)
}

func (r *Recipe) AddStep(s ...string) {
	r.Steps = append(r.Steps, s...)
}
