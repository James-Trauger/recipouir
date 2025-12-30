package com.jamestrauger.recipouir;

import static org.assertj.core.api.Assertions.assertThat;

import java.io.IOException;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.json.JsonTest;
import org.springframework.boot.test.json.JacksonTester;

import com.jamestrauger.recipouir.models.User;

@JsonTest
class UserJsonTest {

    @Autowired
    private JacksonTester<User> json;

    @Test
    void userSerializationTest() throws IOException {
        String username = "asoiaf";
        String firstName = "ned";
        String lastName = "stark";
        // honor
        String password = "$2a$10$SaU8AfqGvUoeLIaJ2W7KY.e3ybJ5RC9mkQxiwFN3tBqu2Jj1vb.XW";

        User user = new User(username, firstName, lastName, password);
        user.setId(47L);

        assertThat(json.write(user)).isStrictlyEqualToJson("expected-user.json");
    }

    @Test
    void userDeserializationTest() throws IOException {
        String username = "asoiaf";
        String firstName = "ned";
        String lastName = "stark";
        // honor
        String password = "$2a$10$SaU8AfqGvUoeLIaJ2W7KY.e3ybJ5RC9mkQxiwFN3tBqu2Jj1vb.XW";
        User user = new User(username, firstName, lastName, password);

        user.setId(47L);

        String expected = """
                {
                    "id": 47,
                    "username": "asoiaf",
                    "firstName": "ned",
                    "lastName": "stark"
                }
                """;
        User serializedUser = json.parse(expected).getObject();

        assertThat(serializedUser).isEqualTo(user);
        assertThat(serializedUser.getId()).isEqualTo(47);
        assertThat(serializedUser.getUsername()).isEqualTo(username);
        assertThat(serializedUser.getFirstName()).isEqualTo(firstName);
        assertThat(serializedUser.getLastName()).isEqualTo(lastName);
    }
}
