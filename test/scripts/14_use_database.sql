-- Test USE database/schema handling
-- USE should create schema if not exists and set search_path

USE test_schema;

-- Create a table in the new schema
CREATE TABLE users (
  id INTEGER,
  name VARCHAR
);

INSERT INTO users VALUES (1, 'Alice'), (2, 'Bob');

-- Query the table (should work via search_path)
SELECT * FROM users ORDER BY id;


