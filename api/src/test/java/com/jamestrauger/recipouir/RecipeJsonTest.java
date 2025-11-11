package com.jamestrauger.recipouir;

import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;
import org.junit.jupiter.api.TestInstance.Lifecycle;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.json.JsonTest;
import org.springframework.boot.test.json.JacksonTester;
import org.springframework.boot.test.json.JsonContent;

import com.jamestrauger.recipouir.models.Recipe;
import com.jamestrauger.recipouir.models.Step;
import com.jamestrauger.recipouir.models.User;
import com.jamestrauger.recipouir.models.Ingredient;
import com.jamestrauger.recipouir.models.Fraction;

import static org.assertj.core.api.Assertions.assertThat;

import java.io.IOException;
import java.util.ArrayList;

@JsonTest
@TestInstance(Lifecycle.PER_CLASS)
class RecipeJsonTest {
    
    @Autowired
    private JacksonTester<Recipe> jsonRecipe;
    @Autowired
    private JacksonTester<Step> jsonStep;
    @Autowired
    private JacksonTester<Ingredient> jsonIngredient;

    // sample data
    private User user;
    private Recipe recipe;

    @BeforeAll
    private void instantiateFields() {
        this.user = new User("asoiaf", "ned", "stark");
        this.user.setId(47L);

        this.recipe = new Recipe("Cookies", user);
        this.recipe.setId(99L);
        this.recipe.setIngredients(new ArrayList<Ingredient>());
        this.recipe.setSteps(new ArrayList<Step>());
    }

    @Test
    void stepSerializationTest() throws IOException {
        Step step = new Step(recipe, "mix the ingredients", 1);

        JsonContent<Step> jsonStepObject = jsonStep.write(step);
        assertThat(jsonStepObject)
            .isStrictlyEqualToJson("expected-step.json");
        assertThat(jsonStepObject)
            .extractingJsonPathStringValue("@.description")
            .isEqualTo(step.getDescription());
        assertThat(jsonStepObject)
            .extractingJsonPathValue("@.id")
            .extracting("number")
            .isEqualTo(step.getId().getNumber());
    }

    @Test
    void stepDeserializationTest() throws IOException {
        String expected = """
                {
                    "id": {
                        "number": 1
                    },
                    "description": "mix the ingredients"
                }
                """;
        Step step = new Step(recipe, "mix the ingredients", 1);
        Step parsedStep = jsonStep.parse(expected).getObject();
        // manually set the recipe
        parsedStep.setRecipe(recipe);
        assertThat(parsedStep).isEqualTo(step);
        assertThat(parsedStep.getDescription()).isEqualTo(step.getDescription());
        assertThat(parsedStep.getId()).isEqualTo(step.getId());
    }

    @Test
    void ingredientSerializationTest() throws IOException {
        Ingredient ingredient = new Ingredient(
            "sugar", 
            recipe,
            100,
            new Fraction(1,1),
            "grams"
        );

        JsonContent<Ingredient> jsonIngredientObject = jsonIngredient.write(ingredient);
        assertThat(jsonIngredientObject).isStrictlyEqualToJson("expected-ingredient.json");
        assertThat(jsonIngredientObject)
            .extractingJsonPathValue("@.id")
            .extracting("name")
            .isEqualTo(ingredient.getName());
        assertThat(jsonIngredientObject)
            .extractingJsonPathStringValue("@.unit")
            .isEqualTo(ingredient.getUnit());
        assertThat(jsonIngredientObject)
            .extractingJsonPathNumberValue("@.amount")
            .isEqualTo(ingredient.getAmount());
        assertThat(jsonIngredientObject)
            .extractingJsonPathValue("@.partialAmount")
            .extracting("numerator")
            .isEqualTo(ingredient.getPartialAmount().numerator());
        assertThat(jsonIngredientObject)
            .extractingJsonPathValue("@.partialAmount")
            .extracting("denominator")
            .isEqualTo(ingredient.getPartialAmount().denominator());
    }

    @Test
    void ingredientDeserializationTest() throws IOException {
        String expected = """
                {
                    "id": {
                        "name": "sugar"
                    },
                    "amount": 100,
                    "unit": "grams",
                    "partialAmount": {
                        "numerator": 1,
                        "denominator": 1
                    }
                }
                """;
        // expected
        Ingredient ingredient = new Ingredient(
            "sugar", 
            recipe,
            100,
            new Fraction(1,1),
            "grams"
        );
        //actual
        Ingredient parsedIngredient = jsonIngredient.parse(expected).getObject();
        parsedIngredient.setRecipe(recipe);

        assertThat(parsedIngredient.getId()).isEqualTo(ingredient.getId());
        assertThat(parsedIngredient.getName()).isEqualTo(ingredient.getName());
        assertThat(parsedIngredient.getAmount()).isEqualTo(ingredient.getAmount());
        assertThat(parsedIngredient.getPartialAmount()).isEqualTo(ingredient.getPartialAmount());
        assertThat(parsedIngredient.getUnit()).isEqualTo(ingredient.getUnit());
        assertThat(parsedIngredient).isEqualTo(ingredient);

    }

    @Test
    void recipeSerializationTest() throws IOException {
        // add ingredients
        this.recipe.getIngredients().add(new Ingredient(
            "sugar", 
            recipe,
            100,
            new Fraction(1,1),
            "grams"
        ));
        this.recipe.getIngredients().add(new Ingredient(
            "butter", 
            recipe,
            1,
            new Fraction(1,2),
            "cups"
        ));

        // add steps
        this.recipe.getSteps()
            .add(new Step(recipe, "mix the ingredients", 1));
        this.recipe.getSteps()
            .add(new Step(recipe, "bake at 350", 2));

        JsonContent<Recipe> jsonRecipeObject = jsonRecipe.write(this.recipe);
        assertThat(jsonRecipeObject)
            .isStrictlyEqualToJson("expected-recipe.json");
        assertThat(jsonRecipeObject)
            .extractingJsonPathNumberValue("@.id")
            .isEqualTo(99);
        assertThat(jsonRecipeObject)
            .extractingJsonPathStringValue("@.title")
            .isEqualTo("Cookies");
        assertThat(jsonRecipeObject)
            .extractingJsonPathValue("@.user")
            .extracting("id")
            .isEqualTo(47);
        assertThat(jsonRecipeObject)
            .extractingJsonPathValue("user")
            .extracting("username")
            .isEqualTo(this.user.getUsername());
        assertThat(jsonRecipeObject)
            .extractingJsonPathValue("@.user")
            .extracting("firstName")
            .isEqualTo(this.user.getFirstName());
        assertThat(jsonRecipeObject)
            .extractingJsonPathValue("@.user")
            .extracting("lastName")
            .isEqualTo(this.user.getLastName());
        // assertThat(jsonRecipeObject)
        //     .extractingJsonPathArrayValue("@.ingredients")
        //     .containsOnly(recipe.getIngredients());
    }

    @Test 
    void recipeDeserializationTest() throws IOException {
        String expected = """
                {
                    "id": 99,
                    "title": "Cookies",
                    "ingredients": [
                        {
                            "id": {
                                "name": "sugar"
                            },
                            "amount": 100,
                            "unit": "grams",
                            "partialAmount": {
                                "numerator": 1,
                                "denominator": 1
                            }
                        },
                        {
                            "id": {
                                "name": "butter"
                            },
                            "amount": 1,
                            "unit": "cups",
                            "partialAmount": {
                                "numerator": 1,
                                "denominator": 2
                            }
                        }
                    ],
                    "steps": [
                        {
                            "id": {
                                "number": 1
                            },
                            "description": "mix the ingredients"
                        },
                        {
                            "id": {
                                "number": 2
                            },
                            "description": "bake at 350"
                        }
                    ],
                    "user": {
                        "id": 47,
                        "username": "asoiaf",
                        "firstName": "ned",
                        "lastName": "stark"
                    }
                }
                """;
        
        Recipe serializedRecipe = jsonRecipe.parse(expected).getObject();
        serializedRecipe.setId(99L);

        assertThat(serializedRecipe)
            .isEqualTo(this.recipe);
        assertThat(serializedRecipe.getId())
            .isEqualTo(99);
        assertThat(serializedRecipe.getTitle())
            .isEqualTo("Cookies");
        
        // expected ingredients
        ArrayList<Ingredient> ingredients = new ArrayList<>();
        ingredients.add(new Ingredient(
            "sugar", 
            recipe,
            100,
            new Fraction(1,1),
            "grams"
        ));
        ingredients.add(new Ingredient(
            "butter", 
            recipe,
            1,
            new Fraction(1,2),
            "cups"
        ));
        assertThat(serializedRecipe.getIngredients())
            .containsAll(ingredients);
        
        // expected steps
        ArrayList<Step> steps = new ArrayList<>();
        steps.add(new Step(recipe, "mix the ingredients", 1));
        steps.add(new Step(recipe, "bake at 350", 2));
        assertThat(serializedRecipe.getSteps())
            .containsAll(steps);
    }
}
