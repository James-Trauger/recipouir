package com.jamestrauger.recipouir;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.jdbc.Sql;

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
		System.out.println(response.getBody());
		Number id = documentContext.read("$.id");
		assertThat(id).isEqualTo(99);

		String title = documentContext.read("$.title");
		assertThat(title).isEqualTo("Cookies");
	}

	@Test
	void shouldNotReturnARecipeWithAnUnknownId() {
		ResponseEntity<String> response = restTemplate.getForEntity("/recipes/1000", String.class);

		assertThat(response.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);
		assertThat(response.getBody()).isBlank();
	}

}
