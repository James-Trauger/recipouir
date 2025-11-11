package com.jamestrauger.recipouir.models;

import com.fasterxml.jackson.annotation.JsonIgnore;

import jakarta.persistence.Column;
import jakarta.persistence.Embedded;
import jakarta.persistence.EmbeddedId;
import jakarta.persistence.Entity;
import jakarta.persistence.Table;

@Entity
//@Table(name = "ingredient")
public class Ingredient {

    @EmbeddedId
    private IngredientId id;
    private int amount;
    @Embedded
    private Fraction partialAmount; //fraction
    @Column(length = 15, nullable = false)
    private String unit;

    protected Ingredient() {}
    
    public Ingredient(String name, Recipe recipe, int amount, Fraction partialAmount, String unit) {
        this.id = new IngredientId(name, recipe);
        this.amount = amount;
        this.partialAmount = new Fraction(partialAmount.numerator(), partialAmount.denominator());
        this.unit = unit;
    }
    
    public IngredientId getId() {
        return this.id;
    }

    public int getAmount() {
        return this.amount;
    }

    public Fraction getPartialAmount() {
        return this.partialAmount;
    }

    public String getUnit() {
        return this.unit;
    }

    @JsonIgnore
    public String getName() {
        return this.id.getName();
    }

    // copies the fields of the given id
    public void setId(IngredientId id) {
        this.id = new IngredientId(id.getName(), id.getRecipe());
    }

    public void setAmount(int amount) {
        this.amount = amount;
    }

    public void setPartialAmount(Fraction partialAmount) {
        this.partialAmount = new Fraction(partialAmount.numerator(), partialAmount.denominator());
    }

    public void setUnit(String unit) {
        this.unit = unit;
    }

    public void setRecipe(Recipe recipe) {
        this.id.setRecipe(recipe);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof Ingredient))
            return false;
        
        Ingredient ingredient = (Ingredient) o;
        if (this.id == null || ingredient.id == null)
            return false;
        return this.id.equals(ingredient.getId())
            && this.amount == ingredient.amount 
            && this.partialAmount.equals(ingredient.partialAmount)
            && this.unit.equals(ingredient.unit);
    }

}
