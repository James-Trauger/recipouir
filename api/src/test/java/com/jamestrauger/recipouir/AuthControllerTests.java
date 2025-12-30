package com.jamestrauger.recipouir;

import static org.assertj.core.api.Assertions.assertThat;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.test.context.ActiveProfiles;
import com.jamestrauger.recipouir.models.User;
import com.jamestrauger.recipouir.models.dto.LoginRequest;
import com.jamestrauger.recipouir.models.dto.SignupRequest;
import com.jamestrauger.recipouir.security.JwtUtil;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@ActiveProfiles("test")
public class AuthControllerTests {
    @Autowired
    private TestRestTemplate restTemplate;
    @Autowired
    private JwtUtil jwtUtil;

    @Test
    void shouldCreateANewUser() {

        // signup request
        SignupRequest signupRequest = new SignupRequest("Tyrion", "Lannister", "tylan", "gold");

        ResponseEntity<User> newUserResponse =
                restTemplate.postForEntity("/api/v1/auth/signup", signupRequest, User.class);

        assertThat(newUserResponse.getStatusCode()).isEqualTo(HttpStatus.CREATED);
        User returnedUser = newUserResponse.getBody();

        assertThat(returnedUser).isNotNull();
        assertThat(returnedUser.getFirstName()).isEqualTo(signupRequest.getFirstName());
        assertThat(returnedUser.getLastName()).isEqualTo(signupRequest.getLastName());
        assertThat(returnedUser.getUsername()).isEqualTo(signupRequest.getUsername());
        // no password should be returned
        assertThat(returnedUser.getPassword()).isNull();
        assertThat(newUserResponse.getHeaders().getLocation().toString())
                .isEqualTo("/api/v1/users/" + returnedUser.getUsername());
    }

    @Test
    void userAlreadyExists() {

        // signup request with username that already exists
        SignupRequest signupRequest = new SignupRequest("Tyrion", "Lannister", "asoiaf", "gold");

        ResponseEntity<User> newUserResponse =
                restTemplate.postForEntity("/api/v1/auth/signup", signupRequest, User.class);

        assertThat(newUserResponse.getStatusCode()).isEqualTo(HttpStatus.CONFLICT);
        assertThat(newUserResponse.getBody()).isNull();
    }

    @Test
    void shouldReturnValidJwtToken() {
        LoginRequest login = new LoginRequest("asoiaf", "honor");

        ResponseEntity<String> tokenResponse =
                restTemplate.postForEntity("/api/v1/auth/signin", login, String.class);

        assertThat(tokenResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        String token = tokenResponse.getBody();
        assertThat(token).isNotNull();

        // validate the returned token
        assertThat(jwtUtil.validateJwtToken(token)).isTrue();

        // subject from the user matches the username logged in
        String userFromToken = jwtUtil.getUserFromToken(token);
        assertThat(userFromToken).isEqualTo(login.getUsername());
    }

    @Test
    void shouldNotSignin() {
        LoginRequest login = new LoginRequest("asoiaf", "corrupt");

        ResponseEntity<String> tokenResponse =
                restTemplate.postForEntity("/api/v1/auth/signin", login, String.class);

        assertThat(tokenResponse.getStatusCode()).isEqualTo(HttpStatus.UNAUTHORIZED);
    }
}
