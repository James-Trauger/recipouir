package com.jamestrauger.recipouir.models;

import java.io.Serializable;
import java.util.Objects;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;

import com.jamestrauger.recipouir.models.deserializers.IngredientIdDeserializer;

@Embeddable
@JsonDeserialize(using = IngredientIdDeserializer.class)
public class IngredientId implements Serializable {
    
    @Column(length = 31)
    private String name;
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "recipe_id")
    @JsonIgnore
    private Recipe recipe;

    protected IngredientId() {}

    public IngredientId(String name, Recipe recipe) {
        this.name = name;
        this.recipe = recipe;
    }

    public String getName() {
        return this.name;
    }

    public Recipe getRecipe() {
        return this.recipe;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setRecipe(Recipe recipe) {
        this.recipe = recipe;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof IngredientId))
            return false;
        IngredientId id = (IngredientId) o;

        if (this.recipe == null)
            return false;

        return name.equals(id.getName()) && this.recipe.getId().equals(id.getRecipe().getId());
    }

    @Override
    public int hashCode() {
        //works for object composition      ???  
        return Objects.hash(this.name, this.recipe);
    }
}
