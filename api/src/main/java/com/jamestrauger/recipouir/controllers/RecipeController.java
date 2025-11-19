package com.jamestrauger.recipouir.controllers;

import java.net.URI;
import java.util.Optional;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.util.UriComponentsBuilder;
import com.jamestrauger.recipouir.models.Recipe;
import com.jamestrauger.recipouir.repositories.RecipeRepository;

@RestController
@RequestMapping("/recipes")
class RecipeController {

    private final RecipeRepository recipeRepository;

    private RecipeController(RecipeRepository recipeRepository) {
        this.recipeRepository = recipeRepository;
    }

    @GetMapping
    public Iterable<Recipe> findAllRecipes() {
        return this.recipeRepository.findAll();
    }

    @PostMapping
    private ResponseEntity<Void> createRecipe(@RequestBody Recipe recipeRequest,
            UriComponentsBuilder ucb) {
        Recipe savedRecipe = recipeRepository.save(recipeRequest);
        URI locationOfNewRecipe =
                ucb.path("recipes/{id}").buildAndExpand(savedRecipe.getId()).toUri();
        return ResponseEntity.created(locationOfNewRecipe).build();
    }

    @GetMapping("/{requestedId}")
    private ResponseEntity<Recipe> fingById(@PathVariable Long requestedId) {
        Optional<Recipe> recipeOptional = recipeRepository.findById(requestedId);
        if (recipeOptional.isPresent())
            return ResponseEntity.ok(recipeOptional.get());
        else
            return ResponseEntity.notFound().build();
    }
}
