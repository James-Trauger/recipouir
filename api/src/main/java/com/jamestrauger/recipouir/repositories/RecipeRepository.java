package com.jamestrauger.recipouir.repositories;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.repository.CrudRepository;
import org.springframework.data.repository.PagingAndSortingRepository;
import com.jamestrauger.recipouir.models.Recipe;

public interface RecipeRepository extends CrudRepository<Recipe, Long>,
                PagingAndSortingRepository<Recipe, Long> {
        Page<Recipe> findByUserUsername(Pageable pageable, String username);
}
