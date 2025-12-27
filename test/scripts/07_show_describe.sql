-- Test SHOW and DESCRIBE commands
CREATE TABLE employees (
  id INTEGER,
  name VARCHAR,
  department VARCHAR,
  salary DOUBLE
);

SHOW TABLES;

DESCRIBE employees;

