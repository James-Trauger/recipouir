package com.jamestrauger.recipouir.repositories;

import org.springframework.data.repository.CrudRepository;
import com.jamestrauger.recipouir.models.User;

public interface UserRepository extends CrudRepository<User, Long> {

}
