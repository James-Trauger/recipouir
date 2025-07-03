package com.jamestrauger.recipouir.errors;

public class NotFoundException extends RuntimeException {

  public NotFoundException(Long id) {
    super("Could not find resource " + id);
  }
}