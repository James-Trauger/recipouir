package com.jamestrauger.recipouir.models;

import java.io.Serializable;
import java.util.Objects;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;

import jakarta.persistence.Embeddable;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;

import com.jamestrauger.recipouir.models.deserializers.StepIdDeserializer;

@Embeddable
@JsonDeserialize(using = StepIdDeserializer.class)
public class StepId implements Serializable {
    
    private int number;
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "recipe_id")
    @JsonIgnore
    private Recipe recipe;

    protected StepId() {}

    public StepId(Recipe recipe, int number) {
        this.recipe = recipe;
        this.number = number;
    }

    public int getNumber() {
        return this.number;
    }

    public Recipe getRecipe() {
        return this.recipe;
    }

    public void setNumber(int number) {
        this.number = number;
    }

    public void setRecipe(Recipe recipe) {
        this.recipe = recipe;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof StepId))
            return false;
        StepId id = (StepId) o;
        return this.number == id.getNumber() 
            && this.recipe.equals(id.getRecipe());
    }

    @Override
    public int hashCode() {
        //works for object composition      ???  
        return Objects.hash(this.number, this.recipe);
    }

    @Override
    public String toString() {
        return String.format("[number=%d, recipe_id=%d]", 
            this.number, this.getRecipe().getId());
    }
}
