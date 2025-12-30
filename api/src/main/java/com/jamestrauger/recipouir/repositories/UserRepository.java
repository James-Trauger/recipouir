package com.jamestrauger.recipouir.repositories;

import java.util.Optional;
import org.springframework.data.repository.CrudRepository;
import com.jamestrauger.recipouir.models.User;

public interface UserRepository extends CrudRepository<User, Long> {
    Optional<User> findByUsername(String username);

    boolean existsByUsername(String username);
}
