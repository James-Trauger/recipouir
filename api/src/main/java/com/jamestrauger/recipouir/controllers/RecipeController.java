package com.jamestrauger.recipouir.controllers;

import java.net.URI;
import java.util.List;
import java.util.Optional;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
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
@RequestMapping("/api/v1/recipes")
class RecipeController {

    private final RecipeRepository recipeRepository;

    private RecipeController(RecipeRepository recipeRepository) {
        this.recipeRepository = recipeRepository;
    }

    @PostMapping
    private ResponseEntity<Void> createRecipe(@RequestBody Recipe recipeRequest,
            UriComponentsBuilder ucb) {

        // TODO: replace with user authorization middleware
        if (recipeRequest.getUser() == null) {
            return ResponseEntity.badRequest().build();
        }

        Recipe savedRecipe = recipeRepository.save(recipeRequest);
        URI locationOfNewRecipe =
                ucb.path("/api/v1/recipes/{username}/{id}")
                        .buildAndExpand(savedRecipe.getUser().getUsername(), savedRecipe.getId())
                        .toUri();
        return ResponseEntity.created(locationOfNewRecipe).build();
    }

    @GetMapping("/{username}")
    private ResponseEntity<List<Recipe>> findAll(@PathVariable String username, Pageable pageable) {
        Page<Recipe> page = recipeRepository.findByUserUsername(
                PageRequest.of(
                        pageable.getPageNumber(),
                        pageable.getPageSize(),
                        pageable.getSortOr(Sort.by(Sort.Direction.ASC, "title"))),
                username);
        return ResponseEntity.ok(page.getContent());
    }

    // single recipe from username with requestedId
    @GetMapping("/{username}/{requestedId}")
    private ResponseEntity<Recipe> fingById(@PathVariable String username,
            @PathVariable Long requestedId) {
        // TODO: find by username and id
        Optional<Recipe> recipeOptional = recipeRepository.findById(requestedId);
        if (recipeOptional.isPresent())
            return ResponseEntity.ok(recipeOptional.get());
        else
            return ResponseEntity.notFound().build();
    }
}
