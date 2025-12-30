package com.jamestrauger.recipouir.models.dto;

import lombok.Data;

@Data
public class SignupRequest {

    private String firstName;
    private String lastName;
    private String username;
    private String password;


    public SignupRequest(String firstName, String lastName, String username, String rawPassword) {
        this.firstName = firstName;
        this.lastName = lastName;
        this.username = username;
        this.password = rawPassword;
    }
}
