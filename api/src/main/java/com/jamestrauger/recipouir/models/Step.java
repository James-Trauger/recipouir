package com.jamestrauger.recipouir.models;

import jakarta.persistence.EmbeddedId;
import jakarta.persistence.Entity;
import jakarta.validation.constraints.Size;

@Entity
public class Step {
    
    @EmbeddedId
    StepId id;
    @Size(min = 1, max = 255)
    private String description;

    protected Step() {}

    public Step(Recipe recipe, String description, int number) {
        this.description = description;
        this.id = new StepId(recipe, number);
    }
    
    public StepId getId() {
        return this.id;
    }

    public String getDescription() {
        return this.description;
    }

    public void setId(StepId id) {
        this.id = id;
    }

    public void setRecipe(Recipe recipe) {
        this.id.setRecipe(recipe);
    }

    public void setDescription(String description) {
        this.description = description;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof Step))
            return false;
        Step step = (Step) o;
        if (this.id == null || step.id == null)
            return false;
        return this.description.equals(step.description)
            && this.id.equals(step.getId());
    }

    @Override
    public String toString() {
        return String.format("\n[\n\tid=%s, description=%s\n\t]",
         this.getId().toString(), description);
    }
}
