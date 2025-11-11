package com.jamestrauger.recipouir.repositories;


import org.springframework.data.repository.CrudRepository;

import com.jamestrauger.recipouir.models.Recipe;

public interface RecipeRepository extends CrudRepository<Recipe, Long> {
    
}
