package com.jamestrauger.recipouir;

import static org.assertj.core.api.Assertions.assertThat;

import java.util.List;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.jdbc.Sql;

import com.jamestrauger.recipouir.models.Recipe;
import com.jamestrauger.recipouir.models.User;
import com.jayway.jsonpath.DocumentContext;
import com.jayway.jsonpath.JsonPath;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@ActiveProfiles("test")
@Sql(scripts = "classpath:data.sql", executionPhase = Sql.ExecutionPhase.BEFORE_TEST_CLASS)
class RecipeControllerTests {

	@Autowired
	TestRestTemplate restTemplate;

	@Test
	void shouldReturnARecipeWhenDataIsSaved() {
		ResponseEntity<String> response = 
			restTemplate.getForEntity("/recipes/99", String.class);
		
		assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);

		DocumentContext documentContext = 
			JsonPath.parse(response.getBody());
		

		Number id = documentContext.read("$.id");
		assertThat(id).isEqualTo(99);

		String title = documentContext.read("$.title");
		assertThat(title).isEqualTo("Cookies");

		// user fields
		String username = documentContext.read("$.user.username");
		String firstName = documentContext.read("$.user.firstName");
		String lastName = documentContext.read("$.user.lastName");
		assertThat(username).isEqualTo("asoiaf");
		assertThat(firstName).isEqualTo(firstName);
		assertThat(lastName).isEqualTo(lastName);

		// ingredient fields

		List<Number> amounts = documentContext.read("$.ingredients[*].amount");
		assertThat(amounts).containsExactly(100, 1);

		List<String> units = documentContext.read("$.ingredients[*].unit");
		assertThat(units).containsExactly("grams", "cups");

		List<String> names = documentContext.read("$.ingredients[*].id.name");
		assertThat(names).containsExactly("sugar", "butter");

		List<Number> numerators = documentContext.read("$.ingredients[*].partialAmount.numerator");
		assertThat(numerators).containsExactly(1, 1);
		List<Number> denominators = documentContext.read("$.ingredients[*].partialAmount.denominator");
		assertThat(denominators).containsExactly(1, 2);

		// step fields

		List<String> descriptions = documentContext.read("$.steps[*].description");
		assertThat(descriptions).containsExactly("mix the ingredients", "bake at 350");

		List<Number> numbers = documentContext.read("$.steps[*].id.number");
		assertThat(numbers).containsExactly(1, 2);
	}

	@Test
	void shouldNotReturnARecipeWithAnUnknownId() {
		ResponseEntity<String> response = restTemplate.getForEntity("/recipes/1000", String.class);

		assertThat(response.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);
		assertThat(response.getBody()).isBlank();
	}

}
