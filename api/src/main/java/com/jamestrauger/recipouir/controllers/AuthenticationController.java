package com.jamestrauger.recipouir.controllers;

import java.net.URI;
import java.net.URISyntaxException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import com.jamestrauger.recipouir.models.User;
import com.jamestrauger.recipouir.models.dto.LoginRequest;
import com.jamestrauger.recipouir.models.dto.SignupRequest;
import com.jamestrauger.recipouir.repositories.UserRepository;
import com.jamestrauger.recipouir.security.JwtUtil;

@RestController
@RequestMapping("/api/v1/auth")
public class AuthenticationController {

        private AuthenticationManager authenticationManager;
        private UserRepository userRepository;
        private PasswordEncoder encoder;
        private JwtUtil jwtUtils;

        @Autowired
        public AuthenticationController(
                        AuthenticationManager authenticationManager,
                        UserRepository userRepository,
                        PasswordEncoder encoder,
                        JwtUtil jwtUtils) {
                this.authenticationManager = authenticationManager;
                this.userRepository = userRepository;
                this.encoder = encoder;
                this.jwtUtils = jwtUtils;
        }


        @PostMapping("/signin")
        public ResponseEntity<String> authenticateUser(@RequestBody LoginRequest user) {
                // attempt to authenticate the user
                Authentication authentication = authenticationManager.authenticate(
                                new UsernamePasswordAuthenticationToken(
                                                user.getUsername(),
                                                user.getPassword()));

                final UserDetails userDetails = (UserDetails) authentication.getPrincipal();
                return ResponseEntity.ok(jwtUtils.generateToken(userDetails.getUsername()));
        }

        @PostMapping("/signup")
        public ResponseEntity<User> registerUser(@RequestBody SignupRequest user) {
                if (userRepository.existsByUsername(user.getUsername())) {
                        return ResponseEntity.status(HttpStatus.CONFLICT).build();
                }

                final User newUser = new User(
                                user.getUsername(),
                                user.getFirstName(),
                                user.getLastName(),
                                encoder.encode(user.getPassword()));
                userRepository.save(newUser);
                try {
                        URI location = new URI("/api/v1/users/" + user.getUsername());
                        return ResponseEntity.created(location).body(newUser);
                } catch (URISyntaxException e) {
                        return ResponseEntity.internalServerError().build();
                }
        }
}
