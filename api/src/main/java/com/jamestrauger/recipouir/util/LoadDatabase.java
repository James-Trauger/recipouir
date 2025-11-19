package com.jamestrauger.recipouir.util;

import javax.sql.DataSource;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.CommandLineRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;



@Configuration
class LoadDatabase {

  private static final Logger log = LoggerFactory.getLogger(LoadDatabase.class);

  @Bean
  CommandLineRunner dsCommandLineRunner(DataSource dataSource) {
    return args -> {
      System.out.println(dataSource.getConnection().getMetaData().getURL());
    };
  }
}