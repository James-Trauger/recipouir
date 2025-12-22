package com.jamestrauger.recipouir.controllers;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import com.jamestrauger.recipouir.repositories.UserRepository;

@RestController
@RequestMapping("/login")
public class LoginController {

    private final UserRepository userRepository;

    private LoginController(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

}
