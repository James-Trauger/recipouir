CREATE USER test WITH PASSWORD 'test';

CREATE DATABASE recipe_app_test;

\c recipe_app_test

CREATE SCHEMA test AUTHORIZATION test;

ALTER ROLE test set search_path = test;

REVOKE ALL PRIVILEGES ON DATABASE postgres FROM test;