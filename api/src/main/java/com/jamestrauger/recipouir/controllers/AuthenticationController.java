package com.jamestrauger.recipouir.controllers;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import com.jamestrauger.recipouir.repositories.UserRepository;

@RestController
@RequestMapping("/login")
public class AuthenticationController {

    private final UserRepository userRepository;

    private AuthenticationController(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

}
