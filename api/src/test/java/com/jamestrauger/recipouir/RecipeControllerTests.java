package com.jamestrauger.recipouir;

import static org.assertj.core.api.Assertions.assertThat;
import java.net.URI;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import java.util.ArrayList;
import java.util.List;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;
import org.junit.jupiter.api.TestInstance.Lifecycle;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.test.annotation.DirtiesContext;
import org.springframework.test.context.ActiveProfiles;
import com.jamestrauger.recipouir.models.Fraction;
import com.jamestrauger.recipouir.models.Ingredient;
import com.jamestrauger.recipouir.models.Recipe;
import com.jamestrauger.recipouir.models.Step;
import com.jamestrauger.recipouir.models.User;
import com.jamestrauger.recipouir.security.JwtUtil;
import com.jayway.jsonpath.DocumentContext;
import com.jayway.jsonpath.JsonPath;
import net.minidev.json.JSONArray;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@TestInstance(Lifecycle.PER_CLASS)
@ActiveProfiles("test")
// uncomment when using docker as the test database source
// @Sql(scripts = "classpath:data.sql", executionPhase = Sql.ExecutionPhase.BEFORE_TEST_CLASS)
class RecipeControllerTests {

	@Autowired
	TestRestTemplate restTemplate;
	@Autowired
	AuthenticationManager authenticationManager;
	@Autowired
	JwtUtil jwtUtil;

	// sample user from database
	private User user;
	// jwt token
	private String token;

	@BeforeAll
	private void insantiateToken() {
		// password is "honor"
		user = new User("asoiaf", "ned", "stark",
				"honor");
		user.setId(47L);

		Authentication authentication = authenticationManager.authenticate(
				new UsernamePasswordAuthenticationToken(
						user.getUsername(),
						user.getPassword()));

		final UserDetails userDetails = (UserDetails) authentication.getPrincipal();
		this.token = jwtUtil.generateToken(userDetails.getUsername());
	}

	@Test
	void shouldReturnARecipeWhenDataIsSaved() {
		HttpHeaders headers = new HttpHeaders();
		headers.add("Authorization", "Bearer " + token);
		ResponseEntity<String> response =
				restTemplate.exchange("/api/v1/recipes/" + user.getUsername() + "/99",
						HttpMethod.GET,
						new HttpEntity<>(headers), String.class);

		assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);

		DocumentContext documentContext = JsonPath.parse(response.getBody());


		Number id = documentContext.read("$.id");
		assertThat(id).isEqualTo(99);

		String title = documentContext.read("$.title");
		assertThat(title).isEqualTo("Cookies");

		Number servings = documentContext.read("$.servings");
		assertThat(servings).isEqualTo(5);

		// user fields
		String username = documentContext.read("$.user.username");
		String firstName = documentContext.read("$.user.firstName");
		String lastName = documentContext.read("$.user.lastName");
		assertThat(username).isEqualTo(user.getUsername());
		assertThat(firstName).isEqualTo(user.getFirstName());
		assertThat(lastName).isEqualTo(user.getLastName());

		// ingredient fields

		List<Number> amounts = documentContext.read("$.ingredients[*].amount");
		assertThat(amounts).containsExactly(100, 1);

		List<String> units = documentContext.read("$.ingredients[*].unit");
		assertThat(units).containsExactly("grams", "cups");

		List<String> names = documentContext.read("$.ingredients[*].id.name");
		assertThat(names).containsExactly("sugar", "butter");

		List<Number> numerators = documentContext.read("$.ingredients[*].partialAmount.numerator");
		assertThat(numerators).containsExactly(1, 1);
		List<Number> denominators =
				documentContext.read("$.ingredients[*].partialAmount.denominator");
		assertThat(denominators).containsExactly(1, 2);

		// step fields

		List<String> descriptions = documentContext.read("$.steps[*].description");
		assertThat(descriptions).containsExactly("mix the ingredients", "bake at 350");

		List<Number> numbers = documentContext.read("$.steps[*].id.number");
		assertThat(numbers).containsExactly(1, 2);
	}

	@Test
	void shouldNotReturnARecipeWithAnUnknownId() {
		HttpHeaders headers = new HttpHeaders();
		headers.add("Authorization", "Bearer " + token);

		ResponseEntity<String> response =
				restTemplate.exchange("/api/v1/recipes/" + user.getUsername() + "/999",
						HttpMethod.GET,
						new HttpEntity<>(headers), String.class);

		assertThat(response.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);
		assertThat(response.getBody()).isBlank();
	}

	@Test
	@DirtiesContext
	void shouldCreateANewRecipe() {
		Recipe recipe = new Recipe("brownies", user, 3);

		// add ingredients to recipe
		ArrayList<Ingredient> ingredients = new ArrayList<Ingredient>();
		ingredients.add(new Ingredient("cacao", recipe, 1, new Fraction(2, 3), "tbs"));
		ingredients.add(new Ingredient("flour", recipe, 150, new Fraction(1, 1), "grams"));
		recipe.setIngredients(ingredients);

		// add steps to the recipe
		ArrayList<Step> steps = new ArrayList<>();
		steps.add(new Step(recipe, "combine cacao and flour", 1));
		steps.add(new Step(recipe, "mix vigorously", 2));
		recipe.setSteps(steps);

		HttpHeaders headers = new HttpHeaders();
		headers.add("Authorization", "Bearer " + token);
		HttpEntity<Recipe> recipeRequest = new HttpEntity<Recipe>(recipe, headers);
		ResponseEntity<Void> createResponse =
				restTemplate.exchange("/api/v1/recipes", HttpMethod.POST,
						recipeRequest, Void.class);

		assertThat(createResponse.getStatusCode()).isEqualTo(HttpStatus.CREATED);

		// retrieve the newly created recipe
		URI locationOfNewRecipe = createResponse.getHeaders().getLocation();


		ResponseEntity<String> getResponse =
				restTemplate.exchange(locationOfNewRecipe, HttpMethod.GET,
						new HttpEntity<>(headers), String.class);

		DocumentContext documentContext = JsonPath.parse(getResponse.getBody());

		String title = documentContext.read("$.title");
		assertThat(title).isEqualTo(recipe.getTitle());

		// generated id
		Number id = documentContext.read("$.id");
		assertThat(id).isNotNull();

		Number servings = documentContext.read("$.servings");
		assertThat(servings).isEqualTo(3);

		// ingredient fields

		List<Number> amounts = documentContext.read("$.ingredients[*].amount");
		assertThat(amounts).containsExactly(1, 150);

		List<String> units = documentContext.read("$.ingredients[*].unit");
		assertThat(units).containsExactly("tbs", "grams");

		List<String> names = documentContext.read("$.ingredients[*].id.name");
		assertThat(names).containsExactly("cacao", "flour");

		List<Number> numerators = documentContext.read("$.ingredients[*].partialAmount.numerator");
		assertThat(numerators).containsExactly(2, 1);
		List<Number> denominators =
				documentContext.read("$.ingredients[*].partialAmount.denominator");
		assertThat(denominators).containsExactly(3, 1);

		// step fields

		List<String> descriptions = documentContext.read("$.steps[*].description");
		assertThat(descriptions).containsExactly("combine cacao and flour", "mix vigorously");

		List<Number> numbers = documentContext.read("$.steps[*].id.number");
		assertThat(numbers).containsExactly(1, 2);

		// user fields
		String username = documentContext.read("$.user.username");
		String firstName = documentContext.read("$.user.firstName");
		String lastName = documentContext.read("$.user.lastName");
		assertThat(username).isEqualTo(user.getUsername());
		assertThat(firstName).isEqualTo(user.getFirstName());
		assertThat(lastName).isEqualTo(user.getLastName());
	}


	@Test
	void shouldNotCreateRecipeWithInvalidUser() {

		// recipe withour a user
		Recipe recipe = new Recipe("Apple Pie", null, 1);

		HttpHeaders headers = new HttpHeaders();
		headers.add("Authorization", "Bearer " + token);
		HttpEntity<Recipe> recipeRequest = new HttpEntity<Recipe>(recipe, headers);
		ResponseEntity<Void> createResponse =
				restTemplate.exchange("/api/v1/recipes", HttpMethod.POST,
						recipeRequest, Void.class);

		assertThat(createResponse.getStatusCode()).isEqualTo(HttpStatus.BAD_REQUEST);
	}

	@Test
	void shouldReturnListRecipes() {
		ResponseEntity<String> response =
				restTemplate.getForEntity("/api/v1/recipes/" + user.getUsername(),
						String.class);

		assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);

		DocumentContext documentContext = JsonPath.parse(response.getBody());
		int recipeCount = documentContext.read("$.length()");
		assertThat(recipeCount).isEqualTo(3);

		JSONArray ids = documentContext.read("$[*].id");
		assertThat(ids).containsExactlyInAnyOrder(99, 100, 101);

		JSONArray titles = documentContext.read("$[*].title");
		assertThat(titles).containsExactlyInAnyOrder("Cookies", "Cake", "Apple Pie");
	}

	@Test
	void shouldReturnAPageOfRecipes() {
		ResponseEntity<String> response =
				restTemplate.getForEntity(
						"/api/v1/recipes/" + user.getUsername() + "?page=0&size=1", String.class);

		assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);

		DocumentContext documentContext = JsonPath.parse(response.getBody());
		JSONArray page = documentContext.read("$[*]");
		assertThat(page.size()).isEqualTo(1);
	}
}
